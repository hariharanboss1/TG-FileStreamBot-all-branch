package config

import (
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var ValueOf = &config{}

type allowedUsers []int64

func (au *allowedUsers) Decode(value string) error {
	if value == "" {
		return nil
	}
	ids := strings.Split(string(value), ",")
	for _, id := range ids {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}
		*au = append(*au, idInt)
	}
	return nil
}

type config struct {
	APIID            int64    `envconfig:"API_ID" required:"true"`
	APIHash          string   `envconfig:"API_HASH" required:"true"`
	BotToken         string   `envconfig:"BOT_TOKEN" required:"true"`
	LogChannelID     int64    `envconfig:"LOG_CHANNEL" required:"true"`
	Host             string   `envconfig:"HOST" required:"true"`
	Port             int      `envconfig:"PORT" required:"true"`
	AllowedUsers     []int64  `envconfig:"ALLOWED_USERS"`
	ForceSubChannel  string   `envconfig:"FORCE_SUB_CHANNEL"`
	Dev              bool     `envconfig:"DEV" default:"false"`
	HashLength       int      `envconfig:"HASH_LENGTH" default:"6"`
	UseSessionFile   bool     `envconfig:"USE_SESSION_FILE" default:"true"`
	UserSession      string   `envconfig:"USER_SESSION"`
	UsePublicIP      bool     `envconfig:"USE_PUBLIC_IP" default:"false"`
	MultiTokens      []string
}

var botTokenRegex = regexp.MustCompile(`MULTI\_TOKEN\d+=(.*)`)

func (c *config) loadFromEnvFile(log *zap.Logger) {
	envPath := filepath.Clean("fsb.env")
	log.Sugar().Infof("Trying to load ENV vars from %s", envPath)
	err := godotenv.Load(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Sugar().Errorf("ENV file not found: %s", envPath)
			log.Sugar().Info("Please create fsb.env file")
			log.Sugar().Info("For more info, refer: https://github.com/EverythingSuckz/TG-FileStreamBot/tree/golang#setting-up-things")
			log.Sugar().Info("Please ignore this message if you are hosting it in a service like Heroku or other alternatives.")
		} else {
			log.Fatal("Unknown error while parsing env file.", zap.Error(err))
		}
	}
}

func (c *config) SetFlagsFromConfig(cmd *cobra.Command) {
	cmd.Flags().Int64Var(&c.APIID, "api-id", 0, "Telegram API ID")
	cmd.Flags().StringVar(&c.APIHash, "api-hash", "", "Telegram API Hash")
	cmd.Flags().StringVar(&c.BotToken, "bot-token", "", "Telegram Bot Token")
	cmd.Flags().Int64Var(&c.LogChannelID, "log-channel", 0, "Log Channel ID")
	cmd.Flags().StringVar(&c.Host, "host", "", "Host URL")
	cmd.Flags().IntVar(&c.Port, "port", 0, "Port")
	cmd.Flags().StringVar(&c.ForceSubChannel, "force-sub-channel", "", "Force Subscription Channel Username")
	cmd.Flags().Bool("dev", c.Dev, "Enable development mode")
	cmd.Flags().Int("hash-length", c.HashLength, "Hash length in links")
	cmd.Flags().Bool("use-session-file", c.UseSessionFile, "Use session files")
	cmd.Flags().String("user-session", c.UserSession, "Pyrogram user session")
	cmd.Flags().Bool("use-public-ip", c.UsePublicIP, "Use public IP instead of local IP")
	cmd.Flags().String("multi-token-txt-file", "", "Multi token txt file (Not implemented)")
}

func (c *config) loadConfigFromArgs(log *zap.Logger, cmd *cobra.Command) {
	if c.APIID != 0 {
		os.Setenv("API_ID", strconv.FormatInt(c.APIID, 10))
	}
	if c.APIHash != "" {
		os.Setenv("API_HASH", c.APIHash)
	}
	if c.BotToken != "" {
		os.Setenv("BOT_TOKEN", c.BotToken)
	}
	if c.LogChannelID != 0 {
		os.Setenv("LOG_CHANNEL", strconv.FormatInt(c.LogChannelID, 10))
	}
	if c.Host != "" {
		os.Setenv("HOST", c.Host)
	}
	if c.Port != 0 {
		os.Setenv("PORT", strconv.Itoa(c.Port))
	}
	if c.ForceSubChannel != "" {
		os.Setenv("FORCE_SUB_CHANNEL", c.ForceSubChannel)
	}
	dev, _ := cmd.Flags().GetBool("dev")
	if dev {
		os.Setenv("DEV", strconv.FormatBool(dev))
	}
	hashLength, _ := cmd.Flags().GetInt("hash-length")
	if hashLength != 0 {
		os.Setenv("HASH_LENGTH", strconv.Itoa(hashLength))
	}
	useSessionFile, _ := cmd.Flags().GetBool("use-session-file")
	if useSessionFile {
		os.Setenv("USE_SESSION_FILE", strconv.FormatBool(useSessionFile))
	}
	userSession, _ := cmd.Flags().GetString("user-session")
	if userSession != "" {
		os.Setenv("USER_SESSION", userSession)
	}
	usePublicIP, _ := cmd.Flags().GetBool("use-public-ip")
	if usePublicIP {
		os.Setenv("USE_PUBLIC_IP", strconv.FormatBool(usePublicIP))
	}
	multiTokens, _ := cmd.Flags().GetString("multi-token-txt-file")
	if multiTokens != "" {
		os.Setenv("MULTI_TOKEN_TXT_FILE", multiTokens)
		// TODO: Add support for importing tokens from a separate file
	}
}

func (c *config) setupEnvVars(log *zap.Logger, cmd *cobra.Command) {
	c.loadFromEnvFile(log)
	c.loadConfigFromArgs(log, cmd)
	err := envconfig.Process("", c)
	if err != nil {
		log.Fatal("Error while parsing env variables", zap.Error(err))
	}
	var ipBlocked bool
	ip, err := getIP(c.UsePublicIP)
	if err != nil {
		log.Error("Error while getting IP", zap.Error(err))
		ipBlocked = true
	}
	if c.Host == "" {
		c.Host = "http://" + ip + ":" + strconv.Itoa(c.Port)
		if c.UsePublicIP {
			if ipBlocked {
				log.Sugar().Warn("Can't get public IP, using local IP")
			} else {
				log.Sugar().Warn("You are using a public IP, please be aware of the security risks while exposing your IP to the internet.")
				log.Sugar().Warn("Use 'HOST' variable to set a domain name")
			}
		}
		log.Sugar().Info("HOST not set, automatically set to " + c.Host)
	}
	val := reflect.ValueOf(c).Elem()
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "MULTI_TOKEN") {
			c.MultiTokens = append(c.MultiTokens, botTokenRegex.FindStringSubmatch(env)[1])
		}
	}
	val.FieldByName("MultiTokens").Set(reflect.ValueOf(c.MultiTokens))
}

func Load(log *zap.Logger, cmd *cobra.Command) {
	log = log.Named("Config")
	defer log.Info("Loaded config")
	ValueOf.setupEnvVars(log, cmd)
	ValueOf.LogChannelID = int64(stripInt(log, int(ValueOf.LogChannelID)))
	if ValueOf.HashLength == 0 {
		log.Sugar().Info("HASH_LENGTH can't be 0, defaulting to 6")
		ValueOf.HashLength = 6
	}
	if ValueOf.HashLength > 32 {
		log.Sugar().Info("HASH_LENGTH can't be more than 32, changing to 32")
		ValueOf.HashLength = 32
	}
	if ValueOf.HashLength < 5 {
		log.Sugar().Info("HASH_LENGTH can't be less than 5, defaulting to 6")
		ValueOf.HashLength = 6
	}
}

func getIP(public bool) (string, error) {
	var ip string
	var err error
	if public {
		ip, err = GetPublicIP()
	} else {
		ip, err = getInternalIP()
	}
	if ip == "" {
		ip = "localhost"
	}
	if err != nil {
		return "localhost", err
	}
	return ip, nil
}

// https://stackoverflow.com/a/23558495/15807350
func getInternalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", errors.New("no internet connection")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if !checkIfIpAccessible(string(ip)) {
		return string(ip), errors.New("PORT is blocked by firewall")
	}
	return string(ip), nil
}

func checkIfIpAccessible(ip string) bool {
	conn, err := net.Dial("tcp", ip+":80")
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func stripInt(log *zap.Logger, a int) int {
	strA := strconv.Itoa(abs(a))
	lastDigits := strings.Replace(strA, "100", "", 1)
	result, err := strconv.Atoi(lastDigits)
	if err != nil {
		log.Sugar().Fatalln(err)
		return 0
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
