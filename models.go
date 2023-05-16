package bdpan

import (
	"bdpan/common"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/wxnacy/go-pretty"
	"github.com/wxnacy/go-tools"
)

type _Credential struct {
	AppId      string `json:"app_id,omitempty"`
	Credentail string `json:"credentail,omitempty"`
}

func newCredentail(c Credential) *_Credential {
	res := &_Credential{}
	res.AppId = c.AppId
	res.SetCredentail(c)
	return res
}

func (c _Credential) GetCredentail() (*Credential, error) {
	res := &Credential{}
	err := decryptHexToInterface(c.Credentail, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *_Credential) SetCredentail(cre Credential) error {

	str, err := encryptInterfaceToHex(cre)
	if err != nil {
		return err
	}
	c.Credentail = str
	return nil
}

type Credential struct {
	AppId       string `json:"app_id,omitempty"`
	AppKey      string `json:"app_key,omitempty"`
	SecretKey   string `json:"secret_key,omitempty"`
	SignKey     string `json:"sign_key,omitempty"`
	accessToken *AccessToken
}

func (c *Credential) GetAccessToken() (*AccessToken, error) {
	if c.accessToken != nil {
		// fmt.Println("return AccessToken from field")
		return c.accessToken, nil
	}
	m, err := common.ReadFileToMap(tokenPath)
	if err != nil {
		return nil, err
	}

	tokenStr, ok := m[c.AppId]
	if !ok {
		return nil, errors.New(fmt.Sprintf(
			"AppId: %s AccessToken not found", c.AppId))
	}

	token := &AccessToken{}
	err = decryptHexToInterface(tokenStr.(string), token)
	if err != nil {
		return nil, err
	}
	// fmt.Println("return AccessToken from file")
	c.accessToken = token
	return token, nil
}

type AccessToken struct {
	ExpiresIn        int32  `json:"expires_in,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	RefreshTimestamp int64  `json:"refresh_timestamp,omitempty"`
}

func (t AccessToken) Print() {
	fmt.Printf("AccessToken: %s\n", t.AccessToken)
	fmt.Printf("RefreshToken: %s\n", t.RefreshToken)
	fmt.Printf("RefreshTimestamp: %d\n", t.RefreshTimestamp)
}

type Config struct {
	LoginAppId string `json:"login_app_id,omitempty"`
}

func NewConfig(loginAppId string) *Config {
	return &Config{LoginAppId: loginAppId}
}

// ----------------------------
// SyncModel
// ----------------------------

type SyncMode int

const (
	ModeBackup SyncMode = iota
	ModeSync
)

func NewSyncModel(local, remote string, mode SyncMode) *SyncModel {
	item := &SyncModel{
		Local:      local,
		Remote:     remote,
		CreateTime: time.Now(),
	}
	item.Hash = tools.Md5(item.Remote + item.Local)
	item.ID = item.Hash[0:7]
	return item
}

func GetModels() (m map[string]*SyncModel) {
	err := tools.FileReadForInterface(syncPath, &m)
	if err != nil {
		m = make(map[string]*SyncModel)
	}
	return
}

func MustGetModel(id string) (m *SyncModel) {
	models := GetModels()
	m, exits := models[id]
	if !exits {
		panic(fmt.Errorf("%s 不存在", id))
	}
	return m
}

func PrintSyncModelList() {
	modelSlice := make([]*SyncModel, 0)
	for _, f := range GetModels() {
		modelSlice = append(modelSlice, f)
	}
	slice := SyncModelSlice(modelSlice)
	sort.Sort(slice)
	pretty.PrintList(slice)

}

func DeleteSyncModel(id string) error {
	models := GetModels()
	m, flag := models[id]
	if !flag {
		return fmt.Errorf("%s 不存在", id)
	}
	fmt.Println(m.Desc())
	flag = PromptConfirm("确定删除")
	if !flag {
		return nil
	}
	delete(models, id)
	err := SaveModels(models)
	if err != nil {
		return err
	}
	PrintSyncModelList()
	return nil
}

func SaveModels(m map[string]*SyncModel) error {
	return tools.FileWriteWithInterface(syncPath, m)
}

type SyncModel struct {
	ID           string
	Remote       string
	Local        string
	Hash         string
	Mode         SyncMode
	HasHide      bool
	CreateTime   time.Time
	LastSyncTime time.Time
}

func (s SyncModel) getLogContent() string {
	mode := "同步"
	if s.IsBackup() {
		mode = "备份"
	}
	return fmt.Sprintf("%s ==> %s %s", s.Local, s.Remote, mode)
}

func (s SyncModel) BuildPretty() []pretty.Field {
	var data = make([]pretty.Field, 0)
	data = append(data, pretty.Field{
		Name:       "ID",
		Value:      s.ID,
		IsFillLeft: true})
	data = append(data, pretty.Field{Name: "Local", Value: s.Local})
	data = append(data, pretty.Field{Name: "Remote", Value: s.Remote})
	data = append(data, pretty.Field{Name: "Mode", Value: s.GetMode()})
	data = append(data, pretty.Field{Name: "HasHide", Value: fmt.Sprintf("%v", s.HasHide)})
	data = append(data, pretty.Field{Name: "CreateTime", Value: s.GetCreateTime()})
	data = append(data, pretty.Field{Name: "LastSyncTime", Value: s.GetLastSyncTime()})
	return data
}

func (s SyncModel) Desc() string {
	tpl := `------------- {{.ID}} ----------------
      ID: {{.ID}}
   Local: {{.Local}}
  Remote: {{.Remote}}
    Mode: {{.GetMode}}
 HasHide: {{.HasHide}}
    Hash: {{.Hash}}
   CTime: {{.GetCreateTime}}
   LSTime: {{.GetLastSyncTime}}
`
	return tools.FormatTemplate(tpl, s)
}

func (s SyncModel) IsBackup() bool {
	if s.Mode == ModeBackup {
		return true
	} else {
		return false
	}
}

func (s SyncModel) GetCreateTime() string {
	return tools.FormatTimeDT(s.CreateTime)
}

func (s SyncModel) GetLastSyncTime() string {
	return tools.FormatTimeDT(s.LastSyncTime)
}

func (s SyncModel) GetMode() string {
	switch s.Mode {
	case ModeBackup:
		return "backup"
	default:
		return "sync"
	}
}

func (s *SyncModel) Exec() error {
	if s.IsBackup() {
		return s.Backup()
	}
	return nil
}

func (s *SyncModel) Backup() error {
	uTasker := NewUploadTasker(s.Local, s.Remote)
	uTasker.IsIncludeHide = s.HasHide
	uTasker.IsRecursion = true
	err := uTasker.Exec()
	if err != nil {
		return err
	}
	s.LastSyncTime = time.Now()
	models := GetModels()
	models[s.ID] = s
	return SaveModels(models)
}

type SyncModelSlice []*SyncModel

func (s SyncModelSlice) GetWriter() io.Writer {
	return os.Stdout
}

func (s SyncModelSlice) List() []pretty.Pretty {
	slice := make([]pretty.Pretty, 0)
	for _, v := range s {
		slice = append(slice, v)
	}
	return slice
}

func (s SyncModelSlice) Len() int {
	return len(s)
}

func (s SyncModelSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SyncModelSlice) Less(i, j int) bool { return s[i].CreateTime.Before(s[j].CreateTime) }
