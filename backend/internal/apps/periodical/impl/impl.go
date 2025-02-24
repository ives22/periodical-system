package impl

import (
	"context"
	"log"
	"periodical/conf"
	"periodical/internal/apps/periodical"

	"gorm.io/gorm"
)

type PeriodicalImpl struct {
	db *gorm.DB
}

func NewPeriodicalImpl() *PeriodicalImpl {
	return &PeriodicalImpl{
		db: conf.C().MySQL.GetConn().Debug(),
	}
}

// CreatePeriodical 创建期刊信息
func (i *PeriodicalImpl) CreatePeriodical(ctx context.Context, req *periodical.CreatePeriodicalRequest) (*periodical.Periodical, error) {
	// 根据创建传入的参数整理查询数据库获取对应的类型，然后整理数据，再写入到数据库
	// 判断传入的类型是否存在
	// 将传入的类型信息进行整合，然后查询是否存在

	// 根据传入的ID，查询通用分类中是否存在
	categorizes, err := i.QueryCategorizeByIDs(ctx, req.CategorizeIDs)
	if err != nil {
		return nil, err
	}

	// 生成一个期刊对象
	ins := periodical.NewPeriodical(req)
	ins.Categorizes = categorizes // 将查询到ID信息赋值给期刊，通过gorm自动实现关联表更新

	// 写入数据
	if err = i.db.WithContext(ctx).Create(ins).Error; err != nil {
		log.Printf("Create periodical failed, err:%v", err.Error())
		return nil, err
	}

	// 开启事务写入数据
	//err = i.db.Transaction(func(tx *gorm.DB) error {
	//	// 插入 Periodical 表
	//	if err = tx.Create(ins).Error; err != nil {
	//		log.Printf("创建期刊失败, err:%v", err)
	//		return err
	//	}
	//
	//	// 手动获取最新插入的 ID
	//	var fetchedPeriodical periodical.Periodical
	//	if err = tx.Where("name = ?", ins.Name).First(&fetchedPeriodical).Error; err != nil {
	//		log.Printf("查询期刊失败, err: %v", err)
	//		return err
	//	}
	//	// 插入 关联表 periodical_common_categories
	//	var periodicalCommonCategories []periodical.PeriodicalCommonCategories
	//	for _, v := range categoriesList {
	//		periodicalCommonCategories = append(periodicalCommonCategories, periodical.PeriodicalCommonCategories{
	//			PeriodicalId: fetchedPeriodical.PeriodicalId,
	//			CommonId:     v.CommonId,
	//		})
	//	}
	//	if err = tx.Create(&periodicalCommonCategories).Error; err != nil {
	//		return err
	//	}
	//	return nil
	//})

	//if err != nil {
	//	return err
	//}

	return ins, nil
}

// QueryPeriodical 查询期刊列表
func (i *PeriodicalImpl) QueryPeriodical(ctx context.Context, req *periodical.QueryPeriodicalRequest) (*periodical.ResponsePeriodicalSet, error) {
	//query := i.db.WithContext(ctx).Model(&periodical.Periodical{})
	query := i.db.WithContext(ctx).
		Joins("JOIN periodical_categorize AS pc ON periodicals.periodical_id = pc.periodical_id").
		Joins("JOIN categorizes ON pc.categorize_id = categorizes.categorize_id")

	// 如果HasBaseTypeList有值, 构建IN条件的值
	if len(req.HasBaseTypeList) > 0 {
		var categorizeIDs []int64
		for _, baseType := range req.HasBaseTypeList {
			categorizeIDs = append(categorizeIDs, int64(baseType.BaseId))
		}
		query = query.Where("categorizes.categorize_id IN ?", categorizeIDs)
	}

	// 根据期刊批次查询
	if len(req.PeriodicalBatch) > 0 {
		log.Println("期刊批次：", req.PeriodicalBatch)
		query = query.Where("periodicals.periodical_batch IN ?", req.PeriodicalBatch)
	}

	// 构造其它过滤条件
	if req.Invoice == 1 {
		query = query.Where("periodicals.invoice_receipt = '1'")
	} else if req.Invoice == 2 {
		query = query.Where("periodicals.invoice_receipt = '0'")
	}
	if req.Works == 1 {
		query = query.Where("periodicals.works = '1'")
	} else if req.Works == 2 {
		query = query.Where("periodicals.works = '0'")
	}
	if req.Name != "" {
		query = query.Where("periodicals.name like ?", "%"+req.Name+"%")
	}
	if req.PublishTimeStart != "" && req.PublishTimeEnd != "" {
		//query = query.Where("periodicals.publication_time > ?", req.PublishTime)
		query = query.Where("periodicals.publication_time BETWEEN ? AND ? ", req.PublishTimeStart, req.PublishTimeEnd)
	}
	if req.DomesticNumber != "" {
		query = query.Where("periodicals.domestic_number like ?", "%"+req.DomesticNumber+"%")
	}
	if req.InternationalNumber != "" {
		query = query.Where("periodicals.international_number like ?", "%"+req.InternationalNumber+"%")
	}
	if req.CompetentUnit != "" {
		query = query.Where("periodicals.competent_unit like ?", "%"+req.CompetentUnit+"%")
	}
	if req.Organization != "" {
		query = query.Where("periodicals.organization like ?", "%"+req.Organization+"%")
	}
	if req.ColumnSetting != "" {
		query = query.Where("periodicals.column_setting like ?", "%"+req.ColumnSetting+"%")
	}
	if req.CompositeInfluenceFactor != 0 {
		if req.CompositeInfluenceFactor == -99 {
			query = query.Where("periodicals.composite_influence_factor >= ?", 0)
		} else {
			query = query.Where("periodicals.composite_influence_factor >= ?", req.CompositeInfluenceFactor)
		}
	}

	query = query.Preload("Categorizes") // 使用Preload预加载关联的分类信息
	query = query.Group("periodicals.periodical_id")

	// 提前准备set对象
	set := periodical.NewResponsePeriodicalSet()

	// 查询总数
	countQuery := query.Session(&gorm.Session{}) // 创建一个新会话，避免影响后续查询
	//err := query.Model(&periodical.Periodical{}).Count(&set.Total).Error
	err := countQuery.Model(&periodical.Periodical{}).Count(&set.Total).Error
	if err != nil {
		log.Println("查询总数失败, ", err)
		return nil, err
	}
	// 分野查询数据
	dataQuery := query.Session(&gorm.Session{})
	//err = query.Offset(req.Offset()).Limit(req.PageSize).Find(&set.Items).Error
	err = dataQuery.Offset(req.Offset()).Limit(req.PageSize).Find(&set.Items).Error
	if err != nil {
		return nil, err
	}

	return set, nil
}

// DescribePeriodical 查询期刊详情
func (i *PeriodicalImpl) DescribePeriodical(ctx context.Context, req *periodical.DescribePeriodicalRequest) (*periodical.Periodical, error) {
	query := i.db.WithContext(ctx).Model(&periodical.Periodical{}).Preload("Categorizes")
	ins := periodical.NewPeriodical(periodical.NewCreatePeriodicalRequest())

	err := query.Where("periodical_id = ?", req.PeriodicalID).Find(ins).Error
	if err != nil {
		return nil, err
	}

	return ins, nil
}

// UpdatePeriodical 更新期刊
func (i *PeriodicalImpl) UpdatePeriodical(ctx context.Context, req *periodical.UpdatePeriodicalRequest) (*periodical.Periodical, error) {
	// 查询现有的期刊信息
	ins, err := i.DescribePeriodical(ctx, periodical.NewDescribePeriodicalRequest(req.PeriodicalID))
	if err != nil {
		return nil, err
	}

	// 根据更新类型进行更新
	switch req.UpdateMode {
	case periodical.UPDATE_MODE_PUT:
		// 全量更新
		ins.BasePeriodical = req.BasePeriodical

	case periodical.UPDATE_MODE_PATCH:
		// 部分更新逻辑（如果需要）
		// 这里可以根据需要更新部分字段
		log.Println("部分更新")
	}

	// 更新期刊信息
	if err := i.db.Save(ins).Error; err != nil {
		return nil, err
	}

	// 更新分类信息
	if req.CategorizeIDs != nil {
		// 清空现有的分类
		if err := i.db.Model(ins).Association("Categorizes").Clear(); err != nil {
			return nil, err
		}

		// 添加新的分类到关联表
		var newCategories []periodical.PeriodicalCategorize
		for _, categorizeID := range req.CategorizeIDs {
			newCategories = append(newCategories, periodical.PeriodicalCategorize{
				PeriodicalID: ins.PeriodicalID,
				CategorizeID: int(categorizeID),
			})
		}

		// 使用 Create 方法将新的分类添加到期刊的关联表
		if err := i.db.Model(&periodical.PeriodicalCategorize{}).Create(&newCategories).Error; err != nil {
			return nil, err
		}
	}

	return ins, nil
}

func (i *PeriodicalImpl) DeletePeriodical(ctx context.Context, req *periodical.DeletePeriodicalRequest) error {
	return i.db.Transaction(func(tx *gorm.DB) error {
		// 删除期刊的关联信息
		if err := tx.WithContext(ctx).Where("periodical_id = ?", req.PeriodicalID).Delete(&periodical.PeriodicalCategorize{}).Error; err != nil {
			return err
		}

		// 删除期刊
		if err := tx.WithContext(ctx).Where("periodical_id = ?", req.PeriodicalID).Delete(&periodical.Periodical{}).Error; err != nil {
			return err
		}
		return nil
	})

	//return i.db.WithContext(ctx).Where("periodical_id = ?", req.PeriodicalID).
	//	Model(&periodical.Periodical{}).
	//	Delete(&periodical.Periodical{}).Error
}

// QueryCategorizeByIDs 根据ID列表查询 通用类型 是否存在
func (i *PeriodicalImpl) QueryCategorizeByIDs(ctx context.Context, ids []int64) (categorizes []*periodical.Categorize, err error) {
	if err = i.db.WithContext(ctx).Where("categorize_id IN ?", ids).Find(&categorizes).Error; err != nil {
		return nil, err
	}
	return categorizes, nil
}

// QueryCategorize 查询通用分类列表
func (i *PeriodicalImpl) QueryCategorize(ctx context.Context, req *periodical.QueryCategorizeRequest) (*periodical.ResponseCategorizeSet, error) {
	query := i.db.WithContext(ctx).Where("status = ?", req.Status)

	set := periodical.NewResponseCategorizeSet()
	// 查询总数
	err := query.Model(&periodical.Categorize{}).Count(&set.Total).Error
	if err != nil {
		log.Println("查询总数失败", err)
		return nil, err
	}

	// 分页查询数据
	err = query.Offset(req.Offset()).Limit(req.PageSize).Find(&set.Items).Error
	if err != nil {
		return nil, err
	}

	return set, nil
}

// mergeIdSlices 将多个 int64 的切片整合为一个切片返回
func (i *PeriodicalImpl) mergeIdSlices(slices ...[]int64) []int64 {
	var merged []int64
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	return merged
}
