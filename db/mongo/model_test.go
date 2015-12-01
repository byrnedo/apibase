package mongo
import (
	"testing"
	. "github.com/byrnedo/apibase/logger"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)


func TestMain(m *testing.M){

	InitLog(func(o *LogOptions){ o.Level = InfoLevel})
	Init("mongodb://localhost:27017/test_mongo_model", Trace)

	m.Run()

	c := Conn()
	defer c.Close()

	c.DB("test_mongo_model").DropDatabase()
}

type MyTypeEntity struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	FirstName string
	LastName string
}

func (m MyTypeEntity) Collection() string {
	return "mytest"
}

type MyType struct{
	FirstName string
	LastName string
}
func (m MyType) MapToEntity() MongoEntity {
	return MyTypeEntity{
		ID : bson.NewObjectId(),
		FirstName: m.FirstName,
		LastName: m.LastName,
	}

}

func TestCreate(t *testing.T){

	var m = MongoModel{}

	user, err := m.Create(MyType{
		FirstName: "Test",
		LastName: "User",
	})

	if err != nil {
		t.Error("Failed to insert:" + err.Error())
	}

	t.Logf("%+v\n", user)

	mSess := Conn()
	col := mSess.DB("").C("mytest")
	defer mSess.Close()

	foundUser := MyTypeEntity{}
	err = col.FindId(user.(MyTypeEntity).ID).One(&foundUser)
	if err != nil {
		t.Error("Failed to find:" + err.Error())
	}

	if reflect.DeepEqual(user, foundUser) == false {
		t.Errorf("Did not match\nexpected:%+v\n   found:%+v\n", user, foundUser)
	}
}
