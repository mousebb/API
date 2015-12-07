package products

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/middleware"
)

type PartVideo struct {
	ID             int
	PartID         int
	YouTubeVideoId string
	Type           string
	IsPrimary      bool
	TypeIcon       *url.URL
	VideoType      VideoType
}

type VideoType struct {
	ID   int
	Name string
	Icon string
}

var (
	partVideoStmt = `select pv.video,vt.name,pv.isPrimary, vt.icon from PartVideo as pv
				join videoType vt on pv.vTypeID = vt.vTypeID
				where pv.partID = ?`
	createPartVideo  = `INSERT INTO PartVideo (partID, video, vTypeID, isPrimary) VALUES (?,?,?,?)`
	deletePartVideos = `DELETE FROM PartVideo WHERE partID = ?`
)

func (p *Part) GetVideos(ctx *middleware.APIContext) error {
	redis_key := fmt.Sprintf("part:%d:videos:%s", p.ID, ctx.DataContext.BrandString)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Videos); err != nil {
			return nil
		}
	}

	qry, err := ctx.DB.Prepare(partVideoStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	rows, err := qry.Query(p.ID)
	if err != nil {
		return err
	}

	var videos []PartVideo
	for rows.Next() {
		var v PartVideo
		var icon *string
		err = rows.Scan(
			&v.YouTubeVideoId,
			&v.Type,
			&v.IsPrimary,
			&icon)
		if err != nil {
			continue
		}

		v.TypeIcon, _ = url.Parse(*icon)
		videos = append(videos, v)
	}
	defer rows.Close()

	go redis.Setex(redis_key, p.Videos, redis.CacheTimeout)

	return nil
}

func (p *PartVideo) CreatePartVideo(ctx *middleware.APIContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:videos:%s", p.PartID, ctx.DataContext.BrandString))

	stmt, err := ctx.DB.Prepare(createPartVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.PartID, p.YouTubeVideoId, p.VideoType.ID, p.IsPrimary)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)
	return nil
}

func (p *PartVideo) DeleteByPart(ctx *middleware.APIContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:videos:%s", p.PartID, ctx.DataContext.BrandString))

	stmt, err := ctx.DB.Prepare(deletePartVideos)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.PartID)

	return err
}
