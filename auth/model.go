package auth

// import "github.com/spf13/viper"

// type Credential struct {
	// // Access    Access `yaml:"access" json:"access"`
	// AppID     string `yaml:"app_id" json:"app_id"`
	// AppKey    string `yaml:"app_key" json:"app_key"`
	// SecretKey string `yaml:"secret_key" json:"secret_key"`
	// SignKey   string `yaml:"sign_key" json:"sign_key"`
// }

// type Access struct {
	// AccessToken      string `yaml:"access_token" json:"access_token"`
	// ExpiresIn        int    `yaml:"expires_in" json:"expires_in"`
	// RefreshTimestamp int    `yaml:"refresh_timestamp" json:"refresh_timestamp"`
	// RefreshToken     string `yaml:"refresh_token" json:"refresh_token"`
// }

// func NewCredential() *Credential {
	// // 设置 Viper 使用环境变量
	// viper.SetConfigType("yaml") // 如果需要加载 YAML 配置文件
	// viper.AutomaticEnv()        // 自动读取环境变量

	// // 设置环境变量前缀
	// viper.SetEnvPrefix("BDPAN")
	// viper.BindEnv("app_id")
	// viper.BindEnv("app_key")
	// viper.BindEnv("secret_key")
	// viper.BindEnv("sign_key")

	// // 创建 Credential 实例
	// var cred Credential

	// // 反序列化环境变量到结构体
	// cred.AppID = viper.GetString("app_id")
	// cred.AppKey = viper.GetString("app_key")
	// cred.SecretKey = viper.GetString("secret_key")
	// cred.SignKey = viper.GetString("sign_key")
	// return &cred
// }
