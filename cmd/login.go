package cmd

import (
	"bdpan"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wxnacy/gotool"
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
	appId := argAppId

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

func (l LoginCommand) getVipName(vipType int32) string {
	switch vipType {
	case 0:
		return "普通用户"
	case 1:
		return "普通会员"
	case 2:
		return "超级会员"
	}
	return "未知身份"
}

func (l LoginCommand) Run() error {
	// appId := *l.AppId
	appId := argAppId
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
	fmt.Printf("Hello, %s(%s)\n", user.GetNetdiskName(), l.getVipName(user.GetVipType()))
	pan, err := bdpan.PanInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取网盘信息失败 %s\n", err.Error())
		return err
	}
	fmt.Printf("网盘容量 %s/%s\n", gotool.FormatSize(pan.GetUsed()), gotool.FormatSize(pan.GetTotal()))
	fmt.Printf("网盘总容量 %d", pan.GetTotal())
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
