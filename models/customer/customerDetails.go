package customer

import (
	"strconv"

	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/middleware"
)

var (
	getDealerTypes = `select dt.dealer_type, ` + dealerTypeFields + ` from DealerTypes as dt
			join ApiKeyToBrand as atb on atb.brandID = dt.brandID
			join ApiKey as a on a.id = atb.keyID
			&&(a.api_key = ? && (dt.brandID = ? or 0 = ?))`
	getDealerTiers = `select dtr.ID, ` + dealerTierFields + ` from DealerTiers as dtr
			join ApiKeyToBrand as atb on atb.brandID = dtr.brandID
			join ApiKey as a on a.id = atb.keyID
			&&(a.api_key = ? && (dtr.brandID = ? or 0 = ?))`
	getMapIcons   = `select mi.ID, mi.tier, mi.dealer_type, ` + mapIconFields + ` from MapIcons as mi`
	getMapixCodes = ` select mpx.mCodeID, ` + mapixCodeFields + ` from MapixCode as mpx`
	getSalesReps  = ` select sr.salesRepID, ` + salesRepFields + ` from salesRepresentative as sr`
)

func DealerTypeMap(ctx *middleware.APIContext) (map[int]DealerType, error) {
	typeMap := make(map[int]DealerType)
	var err error
	dTypes, err := GetDealerTypes(ctx)
	if err != nil {
		return typeMap, err
	}
	for _, dType := range dTypes {
		typeMap[dType.Id] = dType
		//set redis
		redis_key := "dealerType:" + strconv.Itoa(dType.Id)
		err = redis.Set(redis_key, dType)
	}
	return typeMap, err
}

func GetDealerTypes(ctx *middleware.APIContext) ([]DealerType, error) {
	var dType DealerType
	var dTypes []DealerType

	stmt, err := ctx.DB.Prepare(getDealerTypes)
	if err != nil {
		return dTypes, err
	}
	defer stmt.Close()
	res, err := stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)
	if err != nil {
		return dTypes, err
	}
	for res.Next() {
		err = res.Scan(
			&dType.Id,
			&dType.Type,
			&dType.Online,
			&dType.Show,
			&dType.Label,
		)
		if err != nil {
			return dTypes, err
		}
		dTypes = append(dTypes, dType)
	}
	defer res.Close()
	return dTypes, err
}

func DealerTierMap(ctx *middleware.APIContext) (map[int]DealerTier, error) {
	tierMap := make(map[int]DealerTier)
	var err error
	dTiers, err := GetDealerTiers(ctx)
	if err != nil {
		return tierMap, err
	}
	for _, dTier := range dTiers {
		tierMap[dTier.Id] = dTier
		//set redis
		redis_key := "dealerTier:" + strconv.Itoa(dTier.Id)
		err = redis.Set(redis_key, dTier)
	}
	return tierMap, err
}

func GetDealerTiers(ctx *middleware.APIContext) ([]DealerTier, error) {
	var dTier DealerTier
	var dTiers []DealerTier

	stmt, err := ctx.DB.Prepare(getDealerTiers)
	if err != nil {
		return dTiers, err
	}
	defer stmt.Close()

	res, err := stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)
	if err != nil {
		return dTiers, err
	}
	for res.Next() {
		err = res.Scan(
			&dTier.Id,
			&dTier.Tier,
			&dTier.Sort,
		)
		if err != nil {
			return dTiers, err
		}

		dTiers = append(dTiers, dTier)
	}
	defer res.Close()
	return dTiers, err
}

func GetMapIcons(ctx *middleware.APIContext) ([]MapIcon, error) {
	var mi MapIcon
	var mis []MapIcon

	stmt, err := ctx.DB.Prepare(getMapIcons)
	if err != nil {
		return mis, err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(
			&mi.Id,
			&mi.TierId,
			&mi.DealerTypeId,
			&mi.MapIcon,
			&mi.MapIconShadow,
		)
		if err != nil {
			return mis, err
		}
		mis = append(mis, mi)
	}
	defer res.Close()
	return mis, err
}

func MapixMap(ctx *middleware.APIContext) (map[int]MapixCode, error) {
	mapixMap := make(map[int]MapixCode)
	mcs, err := GetMapixCodes(ctx)
	if err != nil {
		return mapixMap, err
	}
	for _, mc := range mcs {
		mapixMap[mc.ID] = mc
		//set redis
		redis_key := "mapixCode:" + strconv.Itoa(mc.ID)
		err = redis.Set(redis_key, mc)
	}
	return mapixMap, err
}

func GetMapixCodes(ctx *middleware.APIContext) ([]MapixCode, error) {
	var mc MapixCode
	var mcs []MapixCode

	stmt, err := ctx.DB.Prepare(getMapixCodes)
	if err != nil {
		return mcs, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return mcs, err
	}
	for res.Next() {
		err = res.Scan(
			&mc.ID,
			&mc.Code,
			&mc.Description,
		)
		if err != nil {
			return mcs, err
		}
		mcs = append(mcs, mc)
	}
	defer res.Close()
	return mcs, err
}

func SalesRepMap(ctx *middleware.APIContext) (map[int]SalesRepresentative, error) {
	repMap := make(map[int]SalesRepresentative)
	reps, err := GetSalesReps(ctx)
	if err != nil {
		return repMap, err
	}
	for _, rep := range reps {
		repMap[rep.ID] = rep
		//set redis
		redis_key := "salesRep:" + strconv.Itoa(rep.ID)
		err = redis.Set(redis_key, rep)
	}
	return repMap, err
}

func GetSalesReps(ctx *middleware.APIContext) ([]SalesRepresentative, error) {
	var sr SalesRepresentative
	var srs []SalesRepresentative

	stmt, err := ctx.DB.Prepare(getSalesReps)
	if err != nil {
		return srs, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(
			&sr.ID,
			&sr.Name,
			&sr.Code,
		)
		if err != nil {
			return srs, err
		}
		srs = append(srs, sr)
	}
	defer res.Close()
	return srs, err
}
