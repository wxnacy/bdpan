package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wxnacy/bdpan"
	"github.com/wxnacy/go-tools"
)

func runLogin(cmd *cobra.Command, args []string) error {
	return LoginCommand{}.Run()
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录网盘",
	Long:  ``,
	RunE:  runLogin,
}

type LoginCommand struct {
}

func (l LoginCommand) buildCredentail() bdpan.Credential {
	// appId := *l.AppId
	appId := globalArg.AppId

	credential := bdpan.Credential{}
	fmt.Println("请先完善秘钥信息")
	if appId == "" {

		fmt.Print("App Id: ")
		fmt.Scanln(&credential.AppId)
	} else {
		credential.AppId = appId
	}
	fmt.Print("App Key: ")
	fmt.Scanln(&credential.AppKey)
	fmt.Print("Secret Key: ")
	fmt.Scanln(&credential.SecretKey)
	fmt.Print("Sign Key: ")
	fmt.Scanln(&credential.SignKey)
	return credential
}

func (l LoginCommand) Run() error {
	// appId := *l.AppId
	appId := globalArg.AppId
	var err error
	if appId == "" {
		_, err = bdpan.GetCredentails()
	} else {
		_, err = bdpan.GetCredentail(appId)
	}
	// 当查询不到用户时进行添加流程
	if err != nil {
		credential := l.buildCredentail()
		err := bdpan.AddCredentail(credential)
		if err != nil {
			fmt.Fprintf(os.Stderr, "登录失败 %s\n", err.Error())
			return err
		}
		// 获取授权
		err = bdpan.CreateAccessTokenByDeviceCode()
		if err != nil {
			fmt.Fprintf(os.Stderr, "登录失败 %s\n", err.Error())
			return err
		}
	}

	if appId != "" {

		// 设置当前需要使用的 appId
		config, err := bdpan.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "登录失败 %s\n", err.Error())
			return err
		}
		config.LoginAppId = appId
	}

	user, err := bdpan.UserInfo()
	// 获取用户信息失败，可能是授权过期则需要进行授权操作
	if err != nil {
		err = bdpan.CreateAccessTokenByDeviceCode()
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取用户信息失败 %s\n", err.Error())
			return err
		}
	}
	fmt.Printf("Hello, %s(%s)\n", user.GetNetdiskName(), user.GetVipName())
	pan, err := bdpan.PanInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取网盘信息失败 %s\n", err.Error())
		return err
	}
	fmt.Printf("网盘容量 %d(%s/%s)\n", pan.GetUsed(), tools.FormatSize(pan.GetUsed()), tools.FormatSize(pan.GetTotal()))
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
