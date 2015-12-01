package mongo
import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoEntity interface {
	Collection() string
}

type MongoEntityMapper interface {
	MapToEntity() MongoEntity
	Wrap(MongoEntity) interface{}
}

type MongoModel struct {
}

func (m *MongoModel) GetSession(e MongoEntity) (*mgo.Collection, *mgo.Session) {
	mSess := Conn()
	return mSess.DB("").C(e.Collection()), mSess
}

func (m *MongoModel) Create(mEn MongoEntityMapper) (interface{},error){
	var en MongoEntity = mEn.MapToEntity()
	col, ses := m.GetSession(en)
	defer ses.Close()

	return en, col.Insert(en)
}
func (uM *MongoModel) Find(id bson.ObjectId) (interface{}, err error) {
	col, ses := uM.GetSession(en)
	defer ses.Close()

	var u = interface{}
	q := col.FindId(id).One(u)
	return u, q
}
