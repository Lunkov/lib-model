package models

import (
  "flag"
  "testing"
  "io/ioutil"
  "github.com/stretchr/testify/assert"
  "time"
  "reflect"
  "github.com/golang/glog"
  "github.com/google/uuid"
  "github.com/jinzhu/gorm"
  "gopkg.in/yaml.v2"

  "github.com/Lunkov/lib-ref"
  "github.com/Lunkov/lib-auth/base"
  
  "github.com/Lunkov/lib-model/fields"
  "github.com/Lunkov/lib-model/models"
)


////////////////////
// TEST DATA
////////////////////
type DocInfo struct {
  ID             uuid.UUID  `db:"id"            json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`

  CODE           string     `db:"code"          json:"code"          yaml:"code"`  // Номер
  DT             time.Time  `db:"create_date"   json:"create_date"`

  Name           string     `db:"name"          json:"name"          yaml:"name"`  // Наименование
  Description    string     `db:"description"   json:"description"   yaml:"description"`  // Описание
}

func (p DocInfo) DBMigrate(db *gorm.DB, tablename string) {
  db.DropTable(tablename)
  db.Table(tablename).AutoMigrate(&DocInfo{})
}

/////////////////////////
// TESTS
/////////////////////////
func TestCheckModelsAutoMigrate(t *testing.T) {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  flag.Set("v", "9")
  flag.Parse()

  conn := "host=localhost port=15432 user=dbuser dbname=testdb password=password sslmode=disable"
  TEST_MODEL := "test_org"
  
  dbconn := New()
  
  dbconn.BaseAdd("organization", reflect.TypeOf(models.Organization{}))
  dbconn.BaseAdd("doc",          reflect.TypeOf(DocInfo{}))

  dbconn.Init(conn, conn, "./etc.test/")
  
  db := dbconn.GetDBHandleWrite()
  db.DropTable(TEST_MODEL)
  dbconn.DBAutoMigrate(conn)
  
  pc := dbconn.GetClass("aaaaa")
  assert.Equal(t, nil, pc)

  pc = dbconn.GetClass(TEST_MODEL)
  assert.Equal(t, &models.Organization{}, pc)
  

  //userUser  := base.User{EMail: "user@",  Group: "user", Groups: []string{"admin_system", "user_crm"}}
  //user_uid, user_id_ok := aclGetIdUser(&userUser)
  //assert.Equal(t, true, user_id_ok)
  //assert.Equal(t, "11111", user_uid.String())
}

func TestCheckModelsByClass(t *testing.T) {
  flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	// flag.Set("v", "9")
	flag.Parse()

	glog.Info("Logging configured")

  conn := "host=localhost port=15432 user=dbuser dbname=testdb password=password sslmode=disable"
  TEST_MODEL := "test_org"
  
  dbconn := New()
  
  dbconn.BaseAdd("organization", reflect.TypeOf(models.Organization{}))
  dbconn.BaseAdd("doc",          reflect.TypeOf(DocInfo{}))

  dbconn.Init(conn, conn, "./etc.test/")
  
  db := dbconn.GetDBHandleWrite()
  db.DropTable(TEST_MODEL)
  dbconn.DBAutoMigrate(conn)

  /// TEST INSERT WITH CLASS  
  uid1, _ := uuid.Parse("00000002-0003-0004-0005-000000000001")
  test_org1 := models.Organization{ID: uid1, CreatedAt: time.Now(), UpdatedAt: time.Now(), CODE: "test.org.1", Name: "OOO `Organization #1`", AddressLegal: fields.Address{Country: "Russia", Index: "127282", City: "Moscow"}}
  
  // glog.Errorf("!!!>>>>>>>: ORG: test_org1(%v): %v", reflect.TypeOf(test_org1), test_org1)
  db.Table(TEST_MODEL).Create(&test_org1)
  test_org11 := models.Organization{}
  errW1 := db.Table(TEST_MODEL).First(&test_org11, "code = ?", "test.org.1").Error
  assert.Equal(t, nil, errW1)
  // TODO: different time values!!!
  test_org11.CreatedAt = test_org1.CreatedAt
  test_org11.UpdatedAt = test_org1.UpdatedAt
  assert.Equal(t, test_org1, test_org11)
}

func TestCheckModelsByClassUnderMap(t *testing.T) {
  conn := "host=localhost port=15432 user=dbuser dbname=testdb password=password sslmode=disable"
  TEST_MODEL := "test_org"
  
  dbconn := New()
  
  dbconn.BaseAdd("organization", reflect.TypeOf(models.Organization{}))
  dbconn.BaseAdd("doc",          reflect.TypeOf(DocInfo{}))

  dbconn.Init(conn, conn, "./etc.test/")
  
  db := dbconn.GetDBHandleWrite()
  db.DropTable(TEST_MODEL)
  dbconn.DBAutoMigrate(conn)

  /// TEST INSERT WITH CLASS UNDER MAP
  uid2, _ := uuid.Parse("00000002-0003-0004-0005-000000000002")
  test_map_org2 := map[string]interface{}{"id": uid2, "created_at": time.Now(), "updated_at": time.Now(), "code": "test.org.2", "address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "bank.0.account":"111111134583459834573279", "bank.0.bik":"1111111", "bank.1.account":"2111111134583459834573279", "bank.1.bik":"21111111", "name":"OOO `Org #2`"}
  test_org2 := models.Organization{}
  ref.ConvertFromMap(&test_org2, &test_map_org2)
  db.Table(TEST_MODEL).Create(&test_org2)
  test_org21 := models.Organization{}
  errW2 := db.Table(TEST_MODEL).First(&test_org21, "code = ?", "test.org.2").Error
  assert.Equal(t, nil, errW2)
  test_org21.CreatedAt = test_org2.CreatedAt
  test_org21.UpdatedAt = test_org2.UpdatedAt
  assert.Equal(t, test_org2, test_org21)
}

func TestCheckModelsByInterfaceUnderMap(t *testing.T) {
  flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	// flag.Set("v", "9")
	flag.Parse()

	glog.Info("Logging configured")
  
  dbconn := New()
  
  TEST_MODEL := "test_org"
  
  conn := "host=localhost port=15432 user=dbuser dbname=testdb password=password sslmode=disable"
  
  dbconn.BaseAdd("organization", reflect.TypeOf(models.Organization{}))
  dbconn.BaseAdd("doc",          reflect.TypeOf(DocInfo{}))

  dbconn.Init(conn, conn, "./etc.test/")
  
  db := dbconn.GetDBHandleWrite()
  db.DropTable(TEST_MODEL)
  dbconn.DBAutoMigrate(conn)

  /// TEST INSERT WITH INTERFACE UNDER MAP
  uid3, _ := uuid.Parse("00000002-0003-0004-0005-000000000003")
  test_org1 := models.Organization{ID: uid3, CreatedAt: time.Now(), UpdatedAt: time.Now(), CODE: "test.org.3", Name: "OOO `Org #2`", AddressLegal: fields.Address{Country: "Russia", Index: "127282", City: "Moscow"}, Bank: fields.BankAccounts{{BIK: "1111111", Account: "111111134583459834573279"}, {BIK: "21111111", Account: "2111111134583459834573279"}}}
  test_map_org3 := map[string]interface{}{"id": uid3, "created_at": test_org1.CreatedAt, "updated_at": test_org1.UpdatedAt, "code": "test.org.3", "address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "bank.0.account":"111111134583459834573279", "bank.0.bik":"1111111", "bank.1.account":"2111111134583459834573279", "bank.1.bik":"21111111", "name":"OOO `Org #2`"}
  
  test_org3 := dbconn.GetClass(TEST_MODEL)
  ref.ConvertFromMap(test_org3, &test_map_org3)

  db.Table(TEST_MODEL).Create(test_org3)
  
  test_org31 := models.Organization{}
  errW3 := db.Table(TEST_MODEL).First(&test_org31, "code = ?", "test.org.3").Error
  assert.Equal(t, nil, errW3)
  test_org1.CreatedAt = test_org31.CreatedAt
  test_org1.UpdatedAt = test_org31.UpdatedAt
  assert.Equal(t, test_org31, test_org1)
}

func TestCheckModelsByInterface(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	// flag.Set("v", "9")
	flag.Parse()

	glog.Info("Logging configured")
  
  TEST_MODEL := "test_org"
  dbconn := New()
  
  conn := "host=localhost port=15432 user=dbuser dbname=testdb password=password sslmode=disable"
  dbconn.BaseAdd("organization", reflect.TypeOf(models.Organization{}))
  dbconn.BaseAdd("doc",          reflect.TypeOf(DocInfo{}))

  dbconn.Init(conn, conn, "./etc.test/")
  
  db := dbconn.GetDBHandleWrite()
  db.DropTable(TEST_MODEL)
  dbconn.DBAutoMigrate(conn)

  /// TEST INSERT WITH MAP
  uid4, _ := uuid.Parse("00000002-0003-0004-0005-000000000004")

  oki := dbconn.DBInsert(TEST_MODEL, nil, &map[string]interface{}{"id": uid4, "created_at": time.Now(), "updated_at": time.Now(), "code": "test.org.4", "name": "ORG `Test NEW 21`", "address_legal.city": "Moscow", "address_legal.index": "127282", "address_legal.country": "Russia", "bank.0.bik": "0065674747", "bank.0.account": "2342345446560065674747"})
  assert.Equal(t, true, oki)
  
  test_org1 := models.Organization{ID: uid4, CreatedAt: time.Now(), UpdatedAt: time.Now(), CODE: "test.org.4", Name: "ORG `Test NEW 21`", AddressLegal: fields.Address{Country: "Russia", Index: "127282", City: "Moscow"}, Bank: fields.BankAccounts{{BIK: "0065674747", Account: "2342345446560065674747"}}}
  test_org31 := models.Organization{}

  errW3 := db.Table(TEST_MODEL).First(&test_org31, "code = ?", "test.org.4").Error
  assert.Nil(t, errW3)

  test_org31.CreatedAt = test_org1.CreatedAt
  test_org31.UpdatedAt = test_org1.UpdatedAt
  assert.Equal(t, test_org1, test_org31)
}

func TestCheckLoadOrg(t *testing.T) {
  yamlFile, err := ioutil.ReadFile("./etc.test/data/test_org/moslift.msk.yaml")
  assert.Equal(t, nil, err)
  
  orgF := make(map[string]models.Organization)
  err = yaml.Unmarshal(yamlFile, orgF)
  assert.Equal(t, nil, err)
  
  uid1, _ := uuid.Parse("29db360f-045d-11ea-8654-bcaec5b972a6")
  orgFN := models.Organization{ID: uid1,
                        Name: "Мослифт",
                        CODE: "moslift.msk.rf",
                        CEO: fields.CompanyPerson{Position: "Генеральный директор", FirstName: "Вартан", LastName: "Авакян", MiddleName: "Нахапетович"},
                        UrlLogo: "/assets/images/logo/moslift-nosign.png",
                        Description: "ОАО «Мослифт» является крупнейшей в России специализированной организацией, которая осуществляет весь комплекс работ по проектированию, поставке, монтажу и техническому обслуживанию лифтов, инвалидных подъёмников, траволаторов и эскалаторов, автопарковочных и объединенных диспетчерских систем."}
  assert.Equal(t, orgFN, orgF["moslift.msk.rf"])
}

func TestCheckModelsGetList(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	// flag.Set("v", "9")
	flag.Parse()

	glog.Info("Logging configured")
  
  dbconn := New()
  
  dbconn.BaseAdd("organization", reflect.TypeOf(models.Organization{}))
  dbconn.BaseAdd("doc",          reflect.TypeOf(DocInfo{}))

  TEST_MODEL := "test_org"

  assert.Equal(t, 14, dbconn.BaseCount())

  conn := "host=localhost port=15432 user=dbuser dbname=testdb password=password sslmode=disable"
  dbconn.Init(conn, conn, "./etc.test/")

  db := dbconn.GetDBHandleWrite()
  db.DropTable(TEST_MODEL)
  dbconn.DBAutoMigrate(conn)

  dbconn.LoadData()
  
  res, count, ok := dbconn.DBTableGet(TEST_MODEL+"org111", nil, []string{"name", "url_logo"}, []string{"name desc"}, 0, 10)
  
  assert.Equal(t, false, ok)
  assert.Equal(t, 0, count)

  assert.Equal(t, "", string(res))

  res, count, ok = dbconn.DBTableGet(TEST_MODEL, nil, []string{"name", "url_logo"}, []string{"name asc", "url_logo asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 6, count)

  resn := "[{\"name\":\"Vaisala\",\"url_logo\":\"/assets/images/logo/vaisala.png\"},{\"name\":\"Wasteout\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\"},{\"name\":\"АО «Объединенная энергетическая компания»\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\"},{\"name\":\"Департамент ИТ Москвы\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\"},{\"name\":\"Мосводоканал\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\"},{\"name\":\"Мослифт\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\"}]"
  assert.Equal(t, resn, string(res))

  res, count, ok = dbconn.DBTableGet(TEST_MODEL, nil, []string{}, []string{}, 0, 0)
  
  assert.Equal(t, true, ok)
  // resr := "[{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0JDQstCw0LrRj9C9INCS0LDRgNGC0LDQvSDQndCw0YXQsNC/0LXRgtC+0LLQuNGHIiwgInBvc2l0aW9uIjogItCT0LXQvdC10YDQsNC70YzQvdGL0Lkg0LTQuNGA0LXQutGC0L7RgCJ9\",\"chief_accountant\":\"e30=\",\"code\":\"moslift.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.280275Z\",\"deleted_at\":null,\"description\":\"ОАО «Мослифт» является крупнейшей в России специализированной организацией, которая осуществляет весь комплекс работ по проектированию, поставке, монтажу и техническому обслуживанию лифтов, инвалидных подъёмников, траволаторов и эскалаторов, автопарковочных и объединенных диспетчерских систем.\",\"id\":\"YWJjZTljNzUtYjIxYy00NjdmLThkMTItOGY4YTcwNzg3M2Rl\",\"name\":\"Мослифт\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.280275Z\",\"url_icon\":\"\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0J/QntCd0J7QnNCQ0KDQldCd0JrQniDQkNC70LXQutGB0LDQvdC00YAg0JzQuNGF0LDQudC70L7QstC40YciLCAicG9zaXRpb24iOiAi0JPQtdC90LXRgNCw0LvRjNC90YvQuSDQtNC40YDQtdC60YLQvtGAIn0=\",\"chief_accountant\":\"e30=\",\"code\":\"mosvodokanal.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.284174Z\",\"deleted_at\":null,\"description\":\"Акционерное общество \\\\Мосводоканал\\\\ – крупнейшая в России водная компания предоставляющая услуги в сфере водоснабжения и водоотведения около 15 млн. жителей Москвы и Московской области – около 10% всего населения страны.\",\"id\":\"ZjNhMTA1NTQtODVkMi00YWNmLWJmZjItNGNhZWY1OTg5MjA3\",\"name\":\"Мосводоканал\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.284174Z\",\"url_icon\":\"\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0JvRi9GB0LXQvdC60L4g0K3QtNGD0LDRgNC0INCQ0L3QsNGC0L7Qu9GM0LXQstC40YciLCAicG9zaXRpb24iOiAi0JzQuNC90LjRgdGC0YAg0J/RgNCw0LLQuNGC0LXQu9GM0YHRgtCy0LAg0JzQvtGB0LrQstGLIn0=\",\"chief_accountant\":\"e30=\",\"code\":\"it.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.286604Z\",\"deleted_at\":null,\"description\":\"Департамент ИТ Москвы занимается разработкой и реализацией госпрограмм в сфере ИТ, телекоммуникаций и связи, что подразумевает пилотирование и последующее внедрение в отрасли городского хозяйства современных и передовых технологий, которые отвечают актуальным потребностям города и москвичей\",\"id\":\"ZGQ2OWEyNjgtYmIwNi00MTQ5LTgxN2QtMTYxYjM1MGNkMzE4\",\"name\":\"Департамент ИТ Москвы\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.286604Z\",\"url_icon\":\"\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0J/RgNC+0YXQvtGA0L7QsiDQldCy0LPQtdC90LjQuSDQodC10YDQs9C10LXQstC40YciLCAicG9zaXRpb24iOiAi0LPQtdC90LXRgNCw0LvRjNC90YvQuSDQtNC40YDQtdC60YLQvtGAIn0=\",\"chief_accountant\":\"e30=\",\"code\":\"oek.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.2888Z\",\"deleted_at\":null,\"description\":\"АО «Объединенная энергетическая компания» - одна из крупнейших электросетевых компаний Москвы, занимающаяся развитием, эксплуатацией и реконструкцией принадлежащих городу электрических сетей. АО «ОЭК» обеспечивает передачу и распределение электроэнергии, осуществляет технологическое присоединение потребителей, ведет строительство новых электрических сетей.\",\"id\":\"MDRiYWJjYzUtYjBjMi00MWZjLThlMWItYmI2YTEzYmFhMGZm\",\"name\":\"АО «Объединенная энергетическая компания»\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.2888Z\",\"url_icon\":\"\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"e30=\",\"chief_accountant\":\"e30=\",\"code\":\"vaisala.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.291769Z\",\"deleted_at\":null,\"description\":\"Vaisala является мировым лидером в области промышленных измерений и измерений параметров окружающей среды. Опираясь на более чем 79-летний опыт работы, компания Vaisala вносит свой вклад в улучшение качества жизни, предоставляя обширный диапазон инновационных продуктов и услуг по наблюдениям и измерениям для метеорологии, чувствительных к метеорологическим условиям операций и управляемого микроклимата. Компания обслуживает клиентов из более чем 140 стран.\",\"id\":\"ZGFhOTcyYWUtYWI1Ni00YTk0LTg2ODQtZmU1MzY5MTJiNzNl\",\"name\":\"Vaisala\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.291769Z\",\"url_icon\":\"\",\"url_logo\":\"/assets/images/logo/vaisala.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"e30=\",\"chief_accountant\":\"e30=\",\"code\":\"wasteout.rf\",\"created_at\":\"2020-05-27T05:53:55.293807Z\",\"deleted_at\":null,\"description\":\"Оптимизация вывоза твердых коммунальных отходов. Снижение эксплутационных расходов на величину от 20 до 50%.\",\"id\":\"NzFjNzc3NDYtNzE5Ny00NTAwLWEzNjItMjYzZjM1NjU3M2Zj\",\"name\":\"Wasteout\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.293807Z\",\"url_icon\":\"\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\",\"url_main\":\"\"}]"
  // assert.Equal(t, resr, string(res))
  
  // uid4, _ := uuid.Parse("00000002-0003-0004-0005-000000000004")
  ok = dbconn.DBInsert(TEST_MODEL, nil, &map[string]interface{}{"id": "00000002-0003-0004-0005-000000000004", "code": "test.1.org", "name": "Test NEW ORG", "address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "bank.0.account":"111111134583459834573279", "bank.0.bik":"1111111", "bank.1.account":"2111111134583459834573279", "bank.1.bik":"21111111"})
  assert.Equal(t, true, ok)

  ok = dbconn.DBInsert(TEST_MODEL, nil, &map[string]interface{}{"id": "00000002-0003-0004-0005-000000000004", "code": "test.1.org", "name": "Test NEW ORG", "address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "bank.0.account":"111111134583459834573279", "bank.0.bik":"1111111", "bank.1.account":"2111111134583459834573279", "bank.1.bik":"21111111"})
  assert.Equal(t, false, ok)

  resn = ""
  res, ok = dbconn.DBGetItemByCODE(TEST_MODEL, nil, "test.1.org.1111")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))

  resn = "{\"address_legal.city\":\"Moscow\",\"address_legal.country\":\"Russia\",\"address_legal.index\":\"127282\",\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"code\":\"test.1.org\",\"id\":\"00000002-0003-0004-0005-000000000004\",\"name\":\"Test NEW ORG\"}"
  res, ok = dbconn.DBGetItemByCODE(TEST_MODEL, nil, "test.1.org")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  resn = ""
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, nil, "00000002-0003-0004-0005-000000004")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))

  resn = "{\"address_legal.city\":\"Moscow\",\"address_legal.country\":\"Russia\",\"address_legal.index\":\"127282\",\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"code\":\"test.1.org\",\"id\":\"00000002-0003-0004-0005-000000000004\",\"name\":\"Test NEW ORG\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, nil, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))


  ok = dbconn.DBUpdate(TEST_MODEL, nil, &map[string]interface{}{"id": "00000002-0003-0004-0005-000000000004", "code": "test.1.org", "name": "Test NEW ORG UPDATE"})
  assert.Equal(t, true, ok)
  
  resn = "{\"code\":\"test.1.org\",\"id\":\"00000002-0003-0004-0005-000000000004\",\"name\":\"Test NEW ORG UPDATE\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, nil, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  ok = dbconn.DBUpdate(TEST_MODEL, nil, &map[string]interface{}{"id": "00000002-0003-0004-0005-000000000004", "code": "test.1.org", "name": "Test NEW ORG UPDATE", "address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "bank.0.account":"111111134583459834573279", "bank.0.bik":"1111111", "bank.1.account":"2111111134583459834573279", "bank.1.bik":"21111111"})
  assert.Equal(t, true, ok)
  
  resn = "{\"address_legal.city\":\"Moscow\",\"address_legal.country\":\"Russia\",\"address_legal.index\":\"127282\",\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"code\":\"test.1.org\",\"id\":\"00000002-0003-0004-0005-000000000004\",\"name\":\"Test NEW ORG UPDATE\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, nil, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  res, count, ok = dbconn.DBTableGet(TEST_MODEL, nil, []string{"name", "url_logo"}, []string{"name asc", "url_logo asc"}, 0, 10)
  assert.Equal(t, true, ok)
  assert.Equal(t, 7, count)

  resn = "[{\"name\":\"Test NEW ORG UPDATE\"},{\"name\":\"Vaisala\",\"url_logo\":\"/assets/images/logo/vaisala.png\"},{\"name\":\"Wasteout\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\"},{\"name\":\"АО «Объединенная энергетическая компания»\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\"},{\"name\":\"Департамент ИТ Москвы\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\"},{\"name\":\"Мосводоканал\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\"},{\"name\":\"Мослифт\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\"}]"
  assert.Equal(t, resn, string(res))
  
  res, count, ok = dbconn.DBTableGet(TEST_MODEL, nil, []string{"name", "url_logo"}, []string{"name asc", "url_logo asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 7, count)

  resn = "[{\"name\":\"Test NEW ORG UPDATE\"},{\"name\":\"Vaisala\",\"url_logo\":\"/assets/images/logo/vaisala.png\"},{\"name\":\"Wasteout\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\"},{\"name\":\"АО «Объединенная энергетическая компания»\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\"},{\"name\":\"Департамент ИТ Москвы\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\"},{\"name\":\"Мосводоканал\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\"},{\"name\":\"Мослифт\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\"}]"
  assert.Equal(t, resn, string(res))

  res, count, ok = dbconn.DBTableGet(TEST_MODEL, nil, []string{"name", "bank"}, []string{"name asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 7, count)

  resn = "[{\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"name\":\"Test NEW ORG UPDATE\"},{\"name\":\"Vaisala\"},{\"name\":\"Wasteout\"},{\"name\":\"АО «Объединенная энергетическая компания»\"},{\"name\":\"Департамент ИТ Москвы\"},{\"name\":\"Мосводоканал\"},{\"name\":\"Мослифт\"}]"
  assert.Equal(t, resn, string(res))


  res, count, ok = dbconn.DBTableGet(TEST_MODEL, nil, []string{}, []string{"name asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 7, count)

  /// DELETE
  ok = dbconn.DBDeleteItemByID(TEST_MODEL, nil, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)

  resn = ""
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, nil, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))


  res, count, ok = dbconn.DBTableGet(TEST_MODEL, nil, []string{"name"}, []string{"name asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 6, count)

  resn = "[{\"name\":\"Vaisala\"},{\"name\":\"Wasteout\"},{\"name\":\"АО «Объединенная энергетическая компания»\"},{\"name\":\"Департамент ИТ Москвы\"},{\"name\":\"Мосводоканал\"},{\"name\":\"Мослифт\"}]"
  assert.Equal(t, resn, string(res))



  dbconn.Close()
}

func TestCheckModelsGetListComplexUpdate(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	// flag.Set("v", "9")
	flag.Parse()

	glog.Info("Logging configured")
  
  dbconn := New()
  
  userAdmin := base.User{EMail: "admin@", Group: "user_crm", Groups: []string{"admin", "user_crm"}}
  userUser  := base.User{EMail: "user@",  Group: "user", Groups: []string{"admin_system", "user_crm"}}
  userOther := base.User{EMail: "guest@", Group: "user_crm", Groups: []string{"admin_system", "user_crm"}}
  
  TEST_MODEL := "test_org_2"
  
  conn := "host=localhost port=15432 user=dbuser dbname=testdb password=password sslmode=disable"
  dbconn.Init(conn, conn, "./etc.test/")

  assert.Equal(t, 14, dbconn.BaseCount())

  db := dbconn.GetDBHandleWrite()
  db.DropTable(TEST_MODEL)
  dbconn.DBAutoMigrate(conn)

  dbconn.LoadData()
  
  res, count, ok := dbconn.DBTableGet(TEST_MODEL + "org111", &userUser, []string{"name", "url_logo"}, []string{"name desc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 0, count)
  assert.Equal(t, "", string(res))

  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{"name", "url_logo"}, []string{"name asc", "url_logo asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 6, count)

  resn := "[{\"name\":\"Vaisala\",\"url_logo\":\"/assets/images/logo/vaisala.png\"},{\"name\":\"Wasteout\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\"},{\"name\":\"АО «Объединенная энергетическая компания»\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\"},{\"name\":\"Департамент ИТ Москвы\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\"},{\"name\":\"Мосводоканал\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\"},{\"name\":\"Мослифт\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\"}]"
  assert.Equal(t, resn, string(res))
  
  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{}, []string{}, 0, 0)
  
  assert.Equal(t, true, ok)
  // resr := "[{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0JDQstCw0LrRj9C9INCS0LDRgNGC0LDQvSDQndCw0YXQsNC/0LXRgtC+0LLQuNGHIiwgInBvc2l0aW9uIjogItCT0LXQvdC10YDQsNC70YzQvdGL0Lkg0LTQuNGA0LXQutGC0L7RgCJ9\",\"chief_accountant\":\"e30=\",\"code\":\"moslift.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.280275Z\",\"deleted_at\":null,\"description\":\"ОАО «Мослифт» является крупнейшей в России специализированной организацией, которая осуществляет весь комплекс работ по проектированию, поставке, монтажу и техническому обслуживанию лифтов, инвалидных подъёмников, траволаторов и эскалаторов, автопарковочных и объединенных диспетчерских систем.\",\"id\":\"YWJjZTljNzUtYjIxYy00NjdmLThkMTItOGY4YTcwNzg3M2Rl\",\"name\":\"Мослифт\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.280275Z\",\"url_icon\":\"\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0J/QntCd0J7QnNCQ0KDQldCd0JrQniDQkNC70LXQutGB0LDQvdC00YAg0JzQuNGF0LDQudC70L7QstC40YciLCAicG9zaXRpb24iOiAi0JPQtdC90LXRgNCw0LvRjNC90YvQuSDQtNC40YDQtdC60YLQvtGAIn0=\",\"chief_accountant\":\"e30=\",\"code\":\"mosvodokanal.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.284174Z\",\"deleted_at\":null,\"description\":\"Акционерное общество \\\\Мосводоканал\\\\ – крупнейшая в России водная компания предоставляющая услуги в сфере водоснабжения и водоотведения около 15 млн. жителей Москвы и Московской области – около 10% всего населения страны.\",\"id\":\"ZjNhMTA1NTQtODVkMi00YWNmLWJmZjItNGNhZWY1OTg5MjA3\",\"name\":\"Мосводоканал\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.284174Z\",\"url_icon\":\"\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0JvRi9GB0LXQvdC60L4g0K3QtNGD0LDRgNC0INCQ0L3QsNGC0L7Qu9GM0LXQstC40YciLCAicG9zaXRpb24iOiAi0JzQuNC90LjRgdGC0YAg0J/RgNCw0LLQuNGC0LXQu9GM0YHRgtCy0LAg0JzQvtGB0LrQstGLIn0=\",\"chief_accountant\":\"e30=\",\"code\":\"it.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.286604Z\",\"deleted_at\":null,\"description\":\"Департамент ИТ Москвы занимается разработкой и реализацией госпрограмм в сфере ИТ, телекоммуникаций и связи, что подразумевает пилотирование и последующее внедрение в отрасли городского хозяйства современных и передовых технологий, которые отвечают актуальным потребностям города и москвичей\",\"id\":\"ZGQ2OWEyNjgtYmIwNi00MTQ5LTgxN2QtMTYxYjM1MGNkMzE4\",\"name\":\"Департамент ИТ Москвы\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.286604Z\",\"url_icon\":\"\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0J/RgNC+0YXQvtGA0L7QsiDQldCy0LPQtdC90LjQuSDQodC10YDQs9C10LXQstC40YciLCAicG9zaXRpb24iOiAi0LPQtdC90LXRgNCw0LvRjNC90YvQuSDQtNC40YDQtdC60YLQvtGAIn0=\",\"chief_accountant\":\"e30=\",\"code\":\"oek.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.2888Z\",\"deleted_at\":null,\"description\":\"АО «Объединенная энергетическая компания» - одна из крупнейших электросетевых компаний Москвы, занимающаяся развитием, эксплуатацией и реконструкцией принадлежащих городу электрических сетей. АО «ОЭК» обеспечивает передачу и распределение электроэнергии, осуществляет технологическое присоединение потребителей, ведет строительство новых электрических сетей.\",\"id\":\"MDRiYWJjYzUtYjBjMi00MWZjLThlMWItYmI2YTEzYmFhMGZm\",\"name\":\"АО «Объединенная энергетическая компания»\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.2888Z\",\"url_icon\":\"\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"e30=\",\"chief_accountant\":\"e30=\",\"code\":\"vaisala.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.291769Z\",\"deleted_at\":null,\"description\":\"Vaisala является мировым лидером в области промышленных измерений и измерений параметров окружающей среды. Опираясь на более чем 79-летний опыт работы, компания Vaisala вносит свой вклад в улучшение качества жизни, предоставляя обширный диапазон инновационных продуктов и услуг по наблюдениям и измерениям для метеорологии, чувствительных к метеорологическим условиям операций и управляемого микроклимата. Компания обслуживает клиентов из более чем 140 стран.\",\"id\":\"ZGFhOTcyYWUtYWI1Ni00YTk0LTg2ODQtZmU1MzY5MTJiNzNl\",\"name\":\"Vaisala\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.291769Z\",\"url_icon\":\"\",\"url_logo\":\"/assets/images/logo/vaisala.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"e30=\",\"chief_accountant\":\"e30=\",\"code\":\"wasteout.rf\",\"created_at\":\"2020-05-27T05:53:55.293807Z\",\"deleted_at\":null,\"description\":\"Оптимизация вывоза твердых коммунальных отходов. Снижение эксплутационных расходов на величину от 20 до 50%.\",\"id\":\"NzFjNzc3NDYtNzE5Ny00NTAwLWEzNjItMjYzZjM1NjU3M2Zj\",\"name\":\"Wasteout\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.293807Z\",\"url_icon\":\"\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\",\"url_main\":\"\"}]"
  // assert.Equal(t, resr, string(res))
  
  // uid4, _ := uuid.Parse("00000002-0003-0004-0005-000000000004")
  ok = dbconn.DBInsert(TEST_MODEL, &userUser, &map[string]interface{}{"id": "00000002-0003-0004-0005-000000000004", "code": "test.1.org", "name": "Test NEW ORG", "address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "bank.0.account":"111111134583459834573279", "bank.0.bik":"1111111", "bank.1.account":"2111111134583459834573279", "bank.1.bik":"21111111"})

  assert.Equal(t, true, ok)

  ok = dbconn.DBInsert(TEST_MODEL, &userUser, &map[string]interface{}{"id": "00000002-0003-0004-0005-000000000004", "code": "test.1.org", "name": "Test NEW ORG", "address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "bank.0.account":"111111134583459834573279", "bank.0.bik":"1111111", "bank.1.account":"2111111134583459834573279", "bank.1.bik":"21111111"})

  assert.Equal(t, false, ok)

  resn = ""
  res, ok = dbconn.DBGetItemByCODE(TEST_MODEL, &userUser, "test.1.org.1111")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))

  resn = "{\"address_legal.city\":\"Moscow\",\"address_legal.country\":\"Russia\",\"address_legal.index\":\"127282\",\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"code\":\"test.1.org\",\"id\":\"00000002-0003-0004-0005-000000000004\",\"name\":\"Test NEW ORG\"}"
  res, ok = dbconn.DBGetItemByCODE(TEST_MODEL, &userUser, "test.1.org")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  resn = ""
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userUser, "00000002-0003-0004-0005-000000004")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))

  resn = "{\"address_legal.city\":\"Moscow\",\"address_legal.country\":\"Russia\",\"address_legal.index\":\"127282\",\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"code\":\"test.1.org\",\"id\":\"00000002-0003-0004-0005-000000000004\",\"name\":\"Test NEW ORG\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userUser, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))


  ok = dbconn.DBUpdate(TEST_MODEL, &userUser, &map[string]interface{}{"id": "00000002-0003-0004-0005-000000000004", "code": "test.1.org", "name": "Test NEW ORG UPDATE"})
  assert.Equal(t, true, ok)
  
  resn = "{\"address_legal.city\":\"Moscow\",\"address_legal.country\":\"Russia\",\"address_legal.index\":\"127282\",\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"code\":\"test.1.org\",\"id\":\"00000002-0003-0004-0005-000000000004\",\"name\":\"Test NEW ORG UPDATE\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userUser, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  ok = dbconn.DBUpdate(TEST_MODEL, &userUser, &map[string]interface{}{"id": "00000002-0003-0004-0005-000000000004", "code": "test.1.org", "name": "Test NEW ORG UPDATE", "address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "bank.0.account":"111111134583459834573279", "bank.0.bik":"1111111", "bank.1.account":"2111111134583459834573279", "bank.1.bik":"21111111"})
  assert.Equal(t, true, ok)
  
  resn = "{\"address_legal.city\":\"Moscow\",\"address_legal.country\":\"Russia\",\"address_legal.index\":\"127282\",\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"code\":\"test.1.org\",\"id\":\"00000002-0003-0004-0005-000000000004\",\"name\":\"Test NEW ORG UPDATE\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userUser, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{"name", "url_logo"}, []string{"name asc", "url_logo asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 7, count)

  resn = "[{\"name\":\"Test NEW ORG UPDATE\"},{\"name\":\"Vaisala\",\"url_logo\":\"/assets/images/logo/vaisala.png\"},{\"name\":\"Wasteout\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\"},{\"name\":\"АО «Объединенная энергетическая компания»\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\"},{\"name\":\"Департамент ИТ Москвы\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\"},{\"name\":\"Мосводоканал\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\"},{\"name\":\"Мослифт\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\"}]"
  assert.Equal(t, resn, string(res))
  
  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{"name", "url_logo"}, []string{"name asc", "url_logo asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 7, count)

  resn = "[{\"name\":\"Test NEW ORG UPDATE\"},{\"name\":\"Vaisala\",\"url_logo\":\"/assets/images/logo/vaisala.png\"},{\"name\":\"Wasteout\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\"},{\"name\":\"АО «Объединенная энергетическая компания»\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\"},{\"name\":\"Департамент ИТ Москвы\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\"},{\"name\":\"Мосводоканал\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\"},{\"name\":\"Мослифт\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\"}]"
  assert.Equal(t, resn, string(res))

  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{"name", "bank"}, []string{"name asc"}, 0, 10)

  assert.Equal(t, true, ok)
  assert.Equal(t, 7, count)

  resn = "[{\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"name\":\"Test NEW ORG UPDATE\"},{\"name\":\"Vaisala\"},{\"name\":\"Wasteout\"},{\"name\":\"АО «Объединенная энергетическая компания»\"},{\"name\":\"Департамент ИТ Москвы\"},{\"name\":\"Мосводоканал\"},{\"name\":\"Мослифт\"}]"
  assert.Equal(t, resn, string(res))


  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{}, []string{"name asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 7, count)

  /// DELETE
  ok = dbconn.DBDeleteItemByID(TEST_MODEL, &userOther, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)

  ok = dbconn.DBDeleteItemByID(TEST_MODEL, &userUser, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)

  ok = dbconn.DBDeleteItemByID(TEST_MODEL, &userAdmin, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)

  resn = ""
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userAdmin, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))


  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userAdmin, []string{"name"}, []string{"name asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 6, count)


  resn = "[{\"name\":\"Vaisala\"},{\"name\":\"Wasteout\"},{\"name\":\"АО «Объединенная энергетическая компания»\"},{\"name\":\"Департамент ИТ Москвы\"},{\"name\":\"Мосводоканал\"},{\"name\":\"Мослифт\"}]"
  assert.Equal(t, resn, string(res))



  dbconn.Close()
}

func TestCheckModelsOwner(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()

	glog.Info("Logging configured")
  uidAdmin, _ := uuid.Parse("00000012-0003-0004-0005-000000000007")
  uidUser, _ := uuid.Parse("00000012-0003-0004-0005-000000000003")
  uidOther, _ := uuid.Parse("00000012-0003-0004-0005-000000000001")
  userAdmin := base.User{ID: uidAdmin, EMail: "admin@", Group: "user_crm", Groups: []string{"admin", "user_crm"}}
  userUser  := base.User{ID: uidUser,  EMail: "user@",  Group: "user", Groups: []string{"admin_system", "user_crm"}}
  userOther := base.User{ID: uidOther, EMail: "guest@", Group: "user_crm", Groups: []string{"admin_system", "user_crm"}}
  
  TEST_MODEL := "test_project"
  dbconn := New()
  
  conn := "host=localhost port=15432 user=dbuser dbname=testdb password=password sslmode=disable"
  dbconn.Init(conn, conn, "./etc.test/")

  assert.Equal(t, 13, dbconn.BaseCount())
  
  db := dbconn.GetDBHandleWrite()
  db.DropTable(TEST_MODEL)
  dbconn.DBAutoMigrate(conn)

  dbconn.LoadData()
  
  res, count, ok := dbconn.DBTableGet(TEST_MODEL + "org111", &userUser, []string{"name", "url_logo"}, []string{"name desc"}, 0, 10)
  
  assert.Equal(t, false, ok)
  assert.Equal(t, 0, count)
  assert.Equal(t, "", string(res))

  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{"name", "description"}, []string{"name asc"}, 0, 10)
  
  assert.Equal(t, true, ok)
  assert.Equal(t, 0, count)
  resn := ""
  //assert.Equal(t, resn, string(res))
  
  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{}, []string{}, 0, 0)
  
  assert.Equal(t, true, ok)
  // resr := "[{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0JDQstCw0LrRj9C9INCS0LDRgNGC0LDQvSDQndCw0YXQsNC/0LXRgtC+0LLQuNGHIiwgInBvc2l0aW9uIjogItCT0LXQvdC10YDQsNC70YzQvdGL0Lkg0LTQuNGA0LXQutGC0L7RgCJ9\",\"chief_accountant\":\"e30=\",\"code\":\"moslift.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.280275Z\",\"deleted_at\":null,\"description\":\"ОАО «Мослифт» является крупнейшей в России специализированной организацией, которая осуществляет весь комплекс работ по проектированию, поставке, монтажу и техническому обслуживанию лифтов, инвалидных подъёмников, траволаторов и эскалаторов, автопарковочных и объединенных диспетчерских систем.\",\"id\":\"YWJjZTljNzUtYjIxYy00NjdmLThkMTItOGY4YTcwNzg3M2Rl\",\"name\":\"Мослифт\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.280275Z\",\"url_icon\":\"\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0J/QntCd0J7QnNCQ0KDQldCd0JrQniDQkNC70LXQutGB0LDQvdC00YAg0JzQuNGF0LDQudC70L7QstC40YciLCAicG9zaXRpb24iOiAi0JPQtdC90LXRgNCw0LvRjNC90YvQuSDQtNC40YDQtdC60YLQvtGAIn0=\",\"chief_accountant\":\"e30=\",\"code\":\"mosvodokanal.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.284174Z\",\"deleted_at\":null,\"description\":\"Акционерное общество \\\\Мосводоканал\\\\ – крупнейшая в России водная компания предоставляющая услуги в сфере водоснабжения и водоотведения около 15 млн. жителей Москвы и Московской области – около 10% всего населения страны.\",\"id\":\"ZjNhMTA1NTQtODVkMi00YWNmLWJmZjItNGNhZWY1OTg5MjA3\",\"name\":\"Мосводоканал\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.284174Z\",\"url_icon\":\"\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0JvRi9GB0LXQvdC60L4g0K3QtNGD0LDRgNC0INCQ0L3QsNGC0L7Qu9GM0LXQstC40YciLCAicG9zaXRpb24iOiAi0JzQuNC90LjRgdGC0YAg0J/RgNCw0LLQuNGC0LXQu9GM0YHRgtCy0LAg0JzQvtGB0LrQstGLIn0=\",\"chief_accountant\":\"e30=\",\"code\":\"it.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.286604Z\",\"deleted_at\":null,\"description\":\"Департамент ИТ Москвы занимается разработкой и реализацией госпрограмм в сфере ИТ, телекоммуникаций и связи, что подразумевает пилотирование и последующее внедрение в отрасли городского хозяйства современных и передовых технологий, которые отвечают актуальным потребностям города и москвичей\",\"id\":\"ZGQ2OWEyNjgtYmIwNi00MTQ5LTgxN2QtMTYxYjM1MGNkMzE4\",\"name\":\"Департамент ИТ Москвы\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.286604Z\",\"url_icon\":\"\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"eyJmaW8iOiAi0J/RgNC+0YXQvtGA0L7QsiDQldCy0LPQtdC90LjQuSDQodC10YDQs9C10LXQstC40YciLCAicG9zaXRpb24iOiAi0LPQtdC90LXRgNCw0LvRjNC90YvQuSDQtNC40YDQtdC60YLQvtGAIn0=\",\"chief_accountant\":\"e30=\",\"code\":\"oek.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.2888Z\",\"deleted_at\":null,\"description\":\"АО «Объединенная энергетическая компания» - одна из крупнейших электросетевых компаний Москвы, занимающаяся развитием, эксплуатацией и реконструкцией принадлежащих городу электрических сетей. АО «ОЭК» обеспечивает передачу и распределение электроэнергии, осуществляет технологическое присоединение потребителей, ведет строительство новых электрических сетей.\",\"id\":\"MDRiYWJjYzUtYjBjMi00MWZjLThlMWItYmI2YTEzYmFhMGZm\",\"name\":\"АО «Объединенная энергетическая компания»\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.2888Z\",\"url_icon\":\"\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"e30=\",\"chief_accountant\":\"e30=\",\"code\":\"vaisala.msk.rf\",\"created_at\":\"2020-05-27T05:53:55.291769Z\",\"deleted_at\":null,\"description\":\"Vaisala является мировым лидером в области промышленных измерений и измерений параметров окружающей среды. Опираясь на более чем 79-летний опыт работы, компания Vaisala вносит свой вклад в улучшение качества жизни, предоставляя обширный диапазон инновационных продуктов и услуг по наблюдениям и измерениям для метеорологии, чувствительных к метеорологическим условиям операций и управляемого микроклимата. Компания обслуживает клиентов из более чем 140 стран.\",\"id\":\"ZGFhOTcyYWUtYWI1Ni00YTk0LTg2ODQtZmU1MzY5MTJiNzNl\",\"name\":\"Vaisala\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.291769Z\",\"url_icon\":\"\",\"url_logo\":\"/assets/images/logo/vaisala.png\",\"url_main\":\"\"},{\"address_billing\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_legal\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"address_shipping\":\"eyJjaXR5IjogIiIsICJyb29tIjogIiIsICJob3VzZSI6ICIiLCAiaW5kZXgiOiAiIiwgInN0cmVldCI6ICIiLCAiY291bnRyeSI6ICIifQ==\",\"bank\":\"bnVsbA==\",\"ceo\":\"e30=\",\"chief_accountant\":\"e30=\",\"code\":\"wasteout.rf\",\"created_at\":\"2020-05-27T05:53:55.293807Z\",\"deleted_at\":null,\"description\":\"Оптимизация вывоза твердых коммунальных отходов. Снижение эксплутационных расходов на величину от 20 до 50%.\",\"id\":\"NzFjNzc3NDYtNzE5Ny00NTAwLWEzNjItMjYzZjM1NjU3M2Zj\",\"name\":\"Wasteout\",\"register\":\"e30=\",\"signer\":\"e30=\",\"support\":\"e30=\",\"updated_at\":\"2020-05-27T05:53:55.293807Z\",\"url_icon\":\"\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\",\"url_main\":\"\"}]"
  // assert.Equal(t, resr, string(res))
  
  // uid4, _ := uuid.Parse("00000002-0003-0004-0005-000000000004")
  ok = dbconn.DBInsert(TEST_MODEL, &userUser, &map[string]interface{}{"id": "00000012-0003-0004-0005-000000000004", "code": "test.1.project", "name": "Test NEW PROJECT"})
  assert.Equal(t, true, ok)

  ok = dbconn.DBInsert(TEST_MODEL, &userUser, &map[string]interface{}{"id": "00000012-0003-0004-0005-000000000004", "code": "test.1.project", "name": "Test NEW PROJECT TWO"})
  assert.Equal(t, false, ok)

  ok = dbconn.DBInsert(TEST_MODEL, &userAdmin, &map[string]interface{}{"id": "00000011-0003-0004-0005-000000000004", "code": "test.admin.project", "name": "Test NEW PROJECT ADMIN"})
  assert.Equal(t, false, ok)

  ok = dbconn.DBInsert(TEST_MODEL, &userOther, &map[string]interface{}{"id": "00000011-0003-0004-0005-000000000014", "code": "test.other.project", "name": "Test NEW PROJECT OTHER"})
  assert.Equal(t, false, ok)

  resn = ""
  res, ok = dbconn.DBGetItemByCODE(TEST_MODEL, &userUser, "test.1.project.1111")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))

  resn = "{\"code\":\"test.1.project\",\"id\":\"00000012-0003-0004-0005-000000000004\",\"name\":\"Test NEW PROJECT\",\"owner.email\":\"user@\",\"owner.id\":\"00000012-0003-0004-0005-000000000003\"}"
  res, ok = dbconn.DBGetItemByCODE(TEST_MODEL, &userUser, "test.1.project")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  resn = ""
  res, ok = dbconn.DBGetItemByCODE(TEST_MODEL, &userAdmin, "test.1.project")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))

  resn = "{\"code\":\"test.1.project\",\"id\":\"00000012-0003-0004-0005-000000000004\",\"name\":\"Test NEW PROJECT\",\"owner.email\":\"user@\",\"owner.id\":\"00000012-0003-0004-0005-000000000003\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userUser, "00000012-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  resn = ""
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userUser, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))

  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userAdmin, "00000012-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))

  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userOther, "00000012-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))

  ok = dbconn.DBUpdate(TEST_MODEL, &userUser, &map[string]interface{}{"id": "00000012-0003-0004-0005-000000000004", "name": "Test PROJECT UPDATE"})
  assert.Equal(t, true, ok)
  
  resn = "{\"code\":\"test.1.project\",\"id\":\"00000012-0003-0004-0005-000000000004\",\"name\":\"Test PROJECT UPDATE\",\"owner.email\":\"user@\",\"owner.id\":\"00000012-0003-0004-0005-000000000003\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userUser, "00000012-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  ok = dbconn.DBUpdate(TEST_MODEL, &userUser, &map[string]interface{}{"id": "00000012-0003-0004-0005-000000000004", "work_groups.0.id": "00000012-0003-0004-0005-000000010104", "work_groups.0.name": "admin"})
  assert.Equal(t, true, ok)

  resn = "{\"code\":\"test.1.project\",\"id\":\"00000012-0003-0004-0005-000000000004\",\"name\":\"Test PROJECT UPDATE\",\"owner.email\":\"user@\",\"owner.id\":\"00000012-0003-0004-0005-000000000003\",\"work_groups.0.id\":\"00000012-0003-0004-0005-000000010104\",\"work_groups.0.name\":\"admin\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userUser, "00000012-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

  resn = "{\"code\":\"test.1.project\",\"id\":\"00000012-0003-0004-0005-000000000004\",\"name\":\"Test PROJECT UPDATE\",\"owner.email\":\"user@\",\"owner.id\":\"00000012-0003-0004-0005-000000000003\",\"work_groups.0.id\":\"00000012-0003-0004-0005-000000010104\",\"work_groups.0.name\":\"admin\"}"
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userAdmin, "00000012-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)
  assert.Equal(t, resn, string(res))

return
  
  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{"name"}, []string{"name asc"}, 0, 10)
  assert.Equal(t, true, ok)
  assert.Equal(t, 2, count)

  resn = "[{\"name\":\"Test NEW ORG UPDATE\"},{\"name\":\"Vaisala\",\"url_logo\":\"/assets/images/logo/vaisala.png\"},{\"name\":\"Wasteout\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\"},{\"name\":\"АО «Объединенная энергетическая компания»\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\"},{\"name\":\"Департамент ИТ Москвы\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\"},{\"name\":\"Мосводоканал\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\"},{\"name\":\"Мослифт\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\"}]"
  assert.Equal(t, resn, string(res))
  
  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{"name", "url_logo"}, []string{"name asc", "url_logo asc"}, 0, 10)
  assert.Equal(t, true, ok)

  if count != 7 {
    t.Error(
      "For", "dbTableGet ORG Model After INSERT",
      "expected", 7,
      "got", count,
    )
  }

  resn = "[{\"name\":\"Test NEW ORG UPDATE\"},{\"name\":\"Vaisala\",\"url_logo\":\"/assets/images/logo/vaisala.png\"},{\"name\":\"Wasteout\",\"url_logo\":\"https://static.wasteout.ru/resources/img/logo.png\"},{\"name\":\"АО «Объединенная энергетическая компания»\",\"url_logo\":\"http://uneco.ru/sites/all/themes/oek_new_theme/img/logo3.png\"},{\"name\":\"Департамент ИТ Москвы\",\"url_logo\":\"https://www.mos.ru/dit/upload/structure/institutions/icon/dit2x.png\"},{\"name\":\"Мосводоканал\",\"url_logo\":\"/assets/images/logo/i-mvk-logo.png\"},{\"name\":\"Мослифт\",\"url_logo\":\"/assets/images/logo/moslift-nosign.png\"}]"
  assert.Equal(t, resn, string(res))

  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{"name", "bank"}, []string{"name asc"}, 0, 10)
  assert.Equal(t, true, ok)

  if count != 7 {
    t.Error(
      "For", "dbTableGet ORG Model After INSERT",
      "expected", 7,
      "got", count,
    )
  }

  resn = "[{\"bank.0.account\":\"111111134583459834573279\",\"bank.0.bik\":\"1111111\",\"bank.1.account\":\"2111111134583459834573279\",\"bank.1.bik\":\"21111111\",\"name\":\"Test NEW ORG UPDATE\"},{\"name\":\"Vaisala\"},{\"name\":\"Wasteout\"},{\"name\":\"АО «Объединенная энергетическая компания»\"},{\"name\":\"Департамент ИТ Москвы\"},{\"name\":\"Мосводоканал\"},{\"name\":\"Мослифт\"}]"
  assert.Equal(t, resn, string(res))


  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userUser, []string{}, []string{"name asc"}, 0, 10)
  
  if ok != true {
    t.Error(
      "For", "dbTableGet ORG Model After INSERT",
      "expected", true,
      "got", ok,
    )
  }

  if count != 7 {
    t.Error(
      "For", "dbTableGet ORG Model After INSERT",
      "expected", 7,
      "got", count,
    )
  }

  /// DELETE
  ok = dbconn.DBDeleteItemByID(TEST_MODEL, &userOther, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)

  ok = dbconn.DBDeleteItemByID(TEST_MODEL, &userUser, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)

  ok = dbconn.DBDeleteItemByID(TEST_MODEL, &userAdmin, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, true, ok)

  resn = ""
  res, ok = dbconn.DBGetItemByID(TEST_MODEL, &userAdmin, "00000002-0003-0004-0005-000000000004")
  assert.Equal(t, false, ok)
  assert.Equal(t, resn, string(res))


  res, count, ok = dbconn.DBTableGet(TEST_MODEL, &userAdmin, []string{"name"}, []string{"name asc"}, 0, 10)
  
  if ok != true {
    t.Error(
      "For", "dbTableGet ORG Model After DELETE",
      "expected", true,
      "got", ok,
    )
  }

  if count != 6 {
    t.Error(
      "For", "dbTableGet ORG Model After DELETE",
      "expected", 6,
      "got", count,
    )
  }

  resn = "[{\"name\":\"Vaisala\"},{\"name\":\"Wasteout\"},{\"name\":\"АО «Объединенная энергетическая компания»\"},{\"name\":\"Департамент ИТ Москвы\"},{\"name\":\"Мосводоканал\"},{\"name\":\"Мослифт\"}]"
  assert.Equal(t, resn, string(res))

  dbconn.Close()
}
