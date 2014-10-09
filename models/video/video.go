package video

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/models/site"
	"github.com/curt-labs/GoAPI/models/products"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Video struct {
	ID           int
	Title        string
	VideoType    VideoType
	Description  string
	DateAdded    time.Time
	DateModified time.Time
	IsPrimary    bool
	Thumbnail    string
	Channels     Channels
	Files        CdnFiles
	Categories   []Category
	Parts        []products.Part
	WebsiteId    int //TODO
}
type Videos []Video

type Channel struct {
	ID           int
	Type         ChannelType
	Link         string
	EmbedCode    string
	ForiegnID    string
	DateAdded    time.Time
	DateModified time.Time
	Title        string
	Description  string
}

type Channels []Channel

type ChannelType struct {
	ID          int
	Name        string
	Description string
}

type CdnFile struct {
	ID           int
	Type         CdnFileType
	Path         string
	Bucket       string
	ObjectName   string
	FileSize     string
	DateAdded    time.Time
	DateModified time.Time
	LastUploaded string
}

type CdnFiles []CdnFile

type CdnFileType struct {
	ID          int
	MimeType    string
	Title       string
	Description string
}

type VideoType struct {
	ID   int
	Name string
	Icon string
}

//TODO categories should be their own entity
type Category struct {
	ID    int
	Title string
}

const (
	videoFields       = ` v.ID, v.subjectTypeID, v.title, v.description, v.dateAdded, v.dateModified, v.isPrimary, v.thumbnail `
	videoTypeFields   = `  vt.name, vt.icon `
	channelFields     = ` c.ID, c.typeID, c.link, c.embedCode, c.foriegnID, c.dateAdded, c.dateModified, c.title, c.desc `
	channelTypeFields = ` ct.name, ct.description `
	cdnFileFields     = `cf.ID, cf.typeID, cf.path, cf.dateAdded, cf.lastUploaded, cf.bucket, cf.objectName, cf.fileSize `
	cdnFileTypeFields = ` cft.mimeType, cft.title, cft.description `
)

var (
	getVideo         = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID WHERE v.ID = ?`
	getAllVideos     = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID`
	getAllCdnFiles   = `SELECT ` + cdnFileFields + `,` + cdnFileTypeFields + ` FROM CdnFile AS cf LEFT JOIN CdnFileType AS cft ON cft.ID = cf.typeID `
	getAllChannels   = `SELECT ` + channelFields + `, ` + channelTypeFields + ` FROM Channel AS c LEFT JOIN ChannelType AS ct ON ct.ID = c.typeID `
	getAllVideoTypes = `SELECT vt.vTypeID, ` + videoTypeFields + ` FROM videoType AS vt`
	getPartVideos    = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID 
						LEFT JOIN VideoJoin AS vj on vj.videoID = v.ID WHERE vj.partID = ? ORDER BY v.title `
	getVideoChannels = `SELECT ` + channelFields + `, ` + channelTypeFields + ` FROM VideoChannels AS vc 
					LEFT JOIN Channel AS c on c.ID = vc.channelID
					LEFT JOIN ChannelType AS ct ON ct.ID = c.typeID
					WHERE vc.videoID = ?`
	getVideoCdns = `SELECT ` + cdnFileFields + `, ` + cdnFileTypeFields + ` FROM CdnFile AS cf 
						LEFT JOIN CdnFileType AS cft ON cft.ID = cf.typeID 
						LEFT JOIN VideoCdnFiles AS vcf ON vcf.cdnID = cf.ID 
						WHERE vcf.videoID = ? `

	createVideo             = `INSERT INTO VideoNew (subjectTypeID, title, description, dateAdded, dateModified, isPrimary, thumbnail) VALUES(?, ?, ?, ?, ?, ?, ?)`
	updateVideo             = `UPDATE videoNew SET subjectTypeID = ?, title = ?, description = ?, isPrimary = ?, thumbnail = ? WHERE ID = ?`
	deleteVideo             = `DELETE FROM videoNew WHERE ID = ?`
	joinVideoCdn            = `INSERT INTO VideoCdnFiles(cdnID, videoID) VALUES(?,?)`
	joinVideoChannel        = `INSERT INTO VideoChannels( channelID, videoID) VALUES(?,?)`
	joinVideoPart           = `INSERT INTO VideoJoin(videoID, partID, catID, websiteID, isPrimary) VALUES(?,?,0,?,?)`
	joinVideoCategory       = `INSERT INTO VideoJoin(videoID, partID, catID, websiteID, isPrimary) VALUES(?,0,?,?,?)`
	deleteVideoCdnJoin      = `DELETE FROM VideoCdnFiles WHERE videoID = ?`
	deleteVideoChannelJoin  = `DELETE FROM VideoChannels WHERE videoID = ?`
	deleteVideoPartJoin     = `DELETE FROM VideoJoin WHERE videoID = ? AND partID = ?`
	deleteVideoCategoryJoin = `DELETE FROM VideoJoin WHERE videoID = ? AND catID = ?`

	//crud
	getChannel        = "SELECT ID, typeID, link, embedCode, foriegnID, dateAdded, dateModified, title, `desc` FROM Channel WHERE ID = ?"
	createChannel     = "INSERT INTO Channel (typeID, link, embedCode, foriegnID, dateAdded, title, `desc`) VALUES (?,?,?,?,?,?,?)"
	updateChannel     = "UPDATE Channel SET typeID = ?, link = ?, embedCode = ?, foriegnID = ?, title = ?, `desc` = ? WHERE ID = ?"
	deleteChannel     = "DELETE FROM Channel WHERE ID = ?"
	getCdn            = `SELECT ` + cdnFileFields + `, ` + cdnFileTypeFields + `, cf.dateModified FROM CdnFile AS cf LEFT JOIN CdnFileType AS cft ON cft.ID = cf.typeID WHERE cf.ID = ?`
	createCdn         = `INSERT INTO CdnFile (typeID, path, dateAdded, lastUploaded, bucket, objectName, fileSize) VALUES (?,?,?,?,?,?,?)`
	updateCdn         = `UPDATE CdnFile SET typeID = ?, path = ?, lastUploaded = ?, bucket = ?, objectName = ?, fileSize = ? WHERE ID = ?`
	deleteCdn         = `DELETE FROM CdnFile WHERE ID = ?`
	getCdnType        = `SELECT ID, mimeType, title, description FROM CdnFileType WHERE ID = ?`
	createCdnType     = `INSERT INTO CdnFileType (mimeType, title, description) VALUES(?,?,?)`
	updateCdnType     = `UPDATE CdnFileType SET mimeType = ?, title = ?, description = ? WHERE ID = ?`
	deleteCdnType     = `DELETE FROM CdnFileType WHERE ID = ?`
	getVideoType      = `SELECT vTypeID, name, icon FROM VideoType WHERE vTypeID = ?`
	createVideoType   = `INSERT INTO VideoType (name, icon) VALUES (?,?)`
	updateVideoType   = `UPDATE VideoType SET name = ?, icon = ? WHERE vTypeID = ?`
	deleteVideoType   = `DELETE FROM VideoType WHERE vTypeID = ?`
	getChannelType    = `SELECT ID, name, description FROM ChannelType WHERE ID = ?`
	createChannelType = `INSERT INTO ChannelType (name, description) VALUES (?,?)`
	updateChannelType = `UPDATE ChannelType SET name = ?, description = ? WHERE ID = ?`
	deleteChannelType = `DELETE FROM ChannelType WHERE ID = ?`
)

//Base Video
func (v *Video) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(v.ID)

	ch := make(chan Video)

	go populateVideo(row, ch)
	*v = <-ch

	return err
}

func GetAllVideos() (vs Videos, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllVideos)
	if err != nil {
		return vs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return vs, err
	}

	ch := make(chan Videos)
	go populateVideos(rows, ch)
	vs = <-ch
	if len(vs) == 0 {
		err = sql.ErrNoRows
	}

	return vs, err
}

func GetPartVideos(p products.Part) (vs Videos, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPartVideos)
	if err != nil {
		return vs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(p.ID)
	if err != nil {
		return vs, err
	}

	ch := make(chan Videos)
	go populateVideos(rows, ch)
	vs = <-ch
	if len(vs) == 0 {
		err = sql.ErrNoRows
	}

	return vs, err
}

func (v *Video) GetChannels() (chs Channels, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return chs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getVideoChannels)
	if err != nil {
		return chs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(v.ID)
	if err != nil {
		return chs, err
	}

	ch := make(chan Channels)
	go populateChannels(rows, ch)
	chs = <-ch

	if len(chs) == 0 {
		err = sql.ErrNoRows
	}

	return chs, err
}

func (v *Video) GetCdnFiles() (cdns CdnFiles, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cdns, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getVideoCdns)
	if err != nil {
		return cdns, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(v.ID)
	if err != nil {
		return cdns, err
	}
	ch := make(chan CdnFiles)
	go populateCdns(rows, ch)
	cdns = <-ch
	return cdns, err
}

func (v *Video) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	v.DateAdded = time.Now()
	res, err := stmt.Exec(v.VideoType.ID, v.Title, v.Description, v.DateAdded, v.DateModified, v.IsPrimary, v.Thumbnail)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	id, err := res.LastInsertId()
	v.ID = int(id)

	// create joins
	fChan := make(chan int)
	chChan := make(chan int)
	catChan := make(chan int)
	pChan := make(chan int)

	go func() (err error) {
		if len(v.Files) > 0 {
			for _, file := range v.Files {
				err = v.CreateJoinFile(file)
			}
		}
		fChan <- 1
		return err
	}()
	go func() (err error) {
		if len(v.Channels) > 0 {
			for _, channel := range v.Channels {
				err = v.CreateJoinChannel(channel)
			}
		}
		chChan <- 1
		return err
	}()
	go func() (err error) {
		if len(v.Categories) > 0 {
			for _, cat := range v.Categories {
				err = v.CreateJoinCategory(cat)
			}
		}
		catChan <- 1
		return err
	}()
	go func() (err error) {
		if len(v.Parts) > 0 {
			for _, part := range v.Parts {
				err = v.CreateJoinPart(part)
			}
		}
		pChan <- 1
		return err
	}()

	<-fChan
	<-chChan
	<-catChan
	<-pChan

	return err
}

func (v *Video) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.VideoType.ID, v.Title, v.Description, v.IsPrimary, v.Thumbnail, v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// delete and create joins
	fChan := make(chan int)
	chChan := make(chan int)
	catChan := make(chan int)
	pChan := make(chan int)

	go func() (err error) {
		err = v.DeleteJoinFiles()
		if err != nil {
			return err
		}
		if len(v.Files) > 0 {
			for _, file := range v.Files {
				err = v.CreateJoinFile(file)
				if err != nil {
					return err
				}
			}
		}
		fChan <- 1
		return nil
	}()
	go func() (err error) {
		err = v.DeleteJoinChannels()
		if err != nil {
			return err
		}
		if len(v.Channels) > 0 {
			for _, channel := range v.Channels {
				err = v.CreateJoinChannel(channel)
				if err != nil {
					return err
				}
			}

		}
		chChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Categories) > 0 {
			for _, cat := range v.Categories {
				err = v.DeleteJoinCategory(cat)
				if err != nil {
					return err
				}
				err = v.CreateJoinCategory(cat)
				if err != nil {
					return err
				}
			}
		}
		catChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Parts) > 0 {
			for _, part := range v.Parts {
				err = v.DeleteJoinPart(part)
				if err != nil {
					return err
				}
				err = v.CreateJoinPart(part)
				if err != nil {
					return err
				}
			}
		}
		pChan <- 1
		return nil
	}()
	<-fChan
	<-chChan
	<-catChan
	<-pChan

	return err
}

func (v *Video) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	//delete and create joins
	fChan := make(chan int)
	chChan := make(chan int)
	catChan := make(chan int)
	pChan := make(chan int)

	go func() (err error) {
		if len(v.Files) > 0 {
			err = v.DeleteJoinFiles()
			if err != nil {
				return err
			}
		}
		fChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Channels) > 0 {
			err = v.DeleteJoinChannels()
			if err != nil {
				return err
			}
		}
		chChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Categories) > 0 {
			for _, cat := range v.Categories {
				err = v.DeleteJoinCategory(cat)
				if err != nil {
					return err
				}
			}
		}
		catChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Parts) > 0 {
			for _, part := range v.Parts {
				err = v.DeleteJoinPart(part)
				if err != nil {
					return err
				}
			}
		}
		pChan <- 1
		return nil
	}()
	<-fChan
	<-chChan
	<-catChan
	<-pChan

	return err
}

func (v *Video) CreateJoinFile(f CdnFile) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(joinVideoCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(f.ID, v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

//Youtube
func (v *Video) CreateJoinChannel(channel Channel) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(joinVideoChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(channel.ID, v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (v *Video) CreateJoinPart(p products.Part) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(joinVideoPart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID, p.ID, v.WebsiteId, v.IsPrimary)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

//HTML5
func (v *Video) CreateJoinCategory(c Category) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(joinVideoCategory)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID, c.ID, v.WebsiteId, v.IsPrimary)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (v *Video) DeleteJoinFiles() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoCdnJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

//Youtube
func (v *Video) DeleteJoinChannels() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoChannelJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (v *Video) DeleteJoinPart(p products.Part) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoPartJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

//HTML5
func (v *Video) DeleteJoinCategory(c Category) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoCategoryJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (c *Channel) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var link, foreign, title *string
	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.Type.ID,
		&link,
		&c.EmbedCode,
		&foreign,
		&c.DateAdded,
		&c.DateModified,
		&title,
		&c.Description,
	)
	if err != nil {
		return err
	}
	if link != nil {
		c.Link = *link
	}
	if foreign != nil {
		c.ForiegnID = *foreign
	}
	if title != nil {
		c.Title = *title
	}
	return err
}

func GetAllChannels() (cs Channels, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllChannels)
	if err != nil {
		return cs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return cs, err
	}

	ch := make(chan Channels)
	go populateChannels(rows, ch)
	cs = <-ch
	if len(cs) == 0 {
		err = sql.ErrNoRows
	}
	return cs, err
}

func (c *Channel) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()
	c.DateAdded = time.Now()
	res, err := stmt.Exec(c.Type.ID, c.Link, c.EmbedCode, c.ForiegnID, c.DateAdded, c.Title, c.Description)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *Channel) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Type.ID, c.Link, c.EmbedCode, c.ForiegnID, c.Title, c.Description, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *Channel) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *CdnFile) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var last, bucket, object, size, desc *string
	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.Type.ID,
		&c.Path,
		&c.DateAdded,
		&last,
		&bucket,
		&object,
		&size,
		&c.Type.MimeType,
		&c.Type.Title,
		&desc,
		&c.DateModified,
	)
	if err != nil {
		return err
	}
	if last != nil {
		c.LastUploaded = *last
	}
	if bucket != nil {
		c.Bucket = *bucket
	}
	if object != nil {
		c.ObjectName = *object
	}
	if size != nil {
		c.FileSize = *size
	}
	if desc != nil {
		c.Type.Description = *desc
	}
	log.Print(c)
	return err
}

func GetAllCdnFiles() (cs CdnFiles, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllCdnFiles)
	if err != nil {
		return cs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return cs, err
	}

	ch := make(chan CdnFiles)
	go populateCdns(rows, ch)
	cs = <-ch
	if len(cs) == 0 {
		err = sql.ErrNoRows
	}
	return cs, err
}

func (c *CdnFile) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()
	c.DateAdded = time.Now()
	res, err := stmt.Exec(c.Type.ID, c.Path, c.DateAdded, c.LastUploaded, c.Bucket, c.ObjectName, c.FileSize)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *CdnFile) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Type.ID, c.Path, c.LastUploaded, c.Bucket, c.ObjectName, c.FileSize, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *CdnFile) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *CdnFileType) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCdnType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var desc *string
	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.MimeType,
		&c.Title,
		&desc,
	)
	if err != nil {
		return err
	}
	if desc != nil {
		c.Description = *desc
	}
	return err
}

func (c *CdnFileType) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createCdnType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.MimeType, c.Title, c.Description)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *CdnFileType) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateCdnType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.MimeType, c.Title, c.Description, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *CdnFileType) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteCdnType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *VideoType) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getVideoType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.Name,
		&c.Icon,
	)
	if err != nil {
		return err
	}
	return err
}

func GetAllVideoTypes() (vts []VideoType, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vts, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllVideoTypes)
	if err != nil {
		return vts, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return vts, err
	}
	var vt VideoType
	var vName, vIcon *string
	for rows.Next() {
		err = rows.Scan(&vt.ID, &vName, &vIcon)
		if err != nil {
			return vts, err
		}
		if vName != nil {
			vt.Name = *vName
		}
		if vIcon != nil {
			vt.Icon = *vIcon
		}
		vts = append(vts, vt)
	}
	return vts, err
}

func (c *VideoType) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createVideoType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.Name, c.Icon)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *VideoType) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateVideoType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Name, c.Icon, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *VideoType) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *ChannelType) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getChannelType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
	)
	if err != nil {
		return err
	}
	return err
}

func (c *ChannelType) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createChannelType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.Name, c.Description)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *ChannelType) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateChannelType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Name, c.Description, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *ChannelType) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteChannelType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

//Populates a video + type
func populateVideo(row *sql.Row, ch chan Video) {
	var v Video
	var tName, tIcon *string
	err := row.Scan(
		&v.ID,
		&v.VideoType.ID,
		&v.Title,
		&v.Description,
		&v.DateAdded,
		&v.DateModified,
		&v.IsPrimary,
		&v.Thumbnail,
		&tName,
		&tIcon,
	)
	if err != nil {
		ch <- v
		return
	}
	if tName != nil {
		v.VideoType.Name = *tName
	}
	if tIcon != nil {
		v.VideoType.Icon = *tIcon
	}

	ch <- v
	return
}

//Populates video and video type fields
func populateVideos(rows *sql.Rows, ch chan Videos) {
	var v Video
	var vs Videos
	var tName, tIcon *string
	for rows.Next() {
		err := rows.Scan(
			&v.ID,
			&v.VideoType.ID,
			&v.Title,
			&v.Description,
			&v.DateAdded,
			&v.DateModified,
			&v.IsPrimary,
			&v.Thumbnail,
			&tName,
			&tIcon,
		)
		if err != nil {
			ch <- Videos{}
			return
		}
		if tName != nil {
			v.VideoType.Name = *tName
		}
		if tIcon != nil {
			v.VideoType.Icon = *tIcon
		}
		vs = append(vs, v)
	}
	ch <- vs
	return
}

//Populates channels and channel type fields
func populateCdns(rows *sql.Rows, ch chan CdnFiles) {
	var c CdnFile
	var cs CdnFiles
	var last, object, bucket, size, tDesc *string
	for rows.Next() {
		err := rows.Scan(
			&c.ID,
			&c.Type.ID,
			&c.Path,
			&c.DateAdded,
			&last,
			&bucket,
			&object,
			&size,
			&c.Type.MimeType,
			&c.Type.Title,
			&tDesc,
		)
		if err != nil {
			return
		}
		if last != nil {
			c.LastUploaded = *last
		}
		if bucket != nil {
			c.Bucket = *bucket
		}
		if object != nil {
			c.ObjectName = *object
		}
		if size != nil {
			c.FileSize = *size
		}
		if tDesc != nil {
			c.Type.Description = *tDesc
		}
		cs = append(cs, c)
	}
	ch <- cs
	return
}

//populate channels
func populateChannels(rows *sql.Rows, ch chan Channels) {
	var chs Channels
	var c Channel
	var link, foreign, title *string
	for rows.Next() {
		err := rows.Scan(
			&c.ID,
			&c.Type.ID,
			&link,
			&c.EmbedCode,
			&foreign,
			&c.DateAdded,
			&c.DateModified,
			&title,
			&c.Description,
			&c.Type.Name,
			&c.Type.Description,
		)
		if err != nil {
			return
		}
		if link != nil {
			c.Link = *link
		}
		if foreign != nil {
			c.ForiegnID = *foreign
		}
		if title != nil {
			c.Title = *title
		}
		chs = append(chs, c)
	}
	ch <- chs
	return
}
