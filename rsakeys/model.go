package rsakeys

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Clay294/forum/flog"
)

const (
	UnitName = "rsakeys"
)

var globalRSAKeyPairsUUIDs = make([]string, 0, 20)

func GlobalRSAKeyPairsUUIDs() []string {
	return globalRSAKeyPairsUUIDs
}

const (
	RSAKeyPairsStringTableName  = "rsakeys"
	DefaultRSAKeyPairsUUIDsFile = "./rsakeys/rsa_key_pairs_uuids.txt"
)

type RSAKeyPair struct {
	PrivateKeyString string
	PublicKeyString  string
	*RSAKeyPairMeta
}

type RSAKeyPairMeta struct {
	Id        int
	UUID      string
	CreatedAt int64
}

func (rsaKP *RSAKeyPair) TableName() string {
	return RSAKeyPairsStringTableName
}

func InitGlobalRSAKeyPairsUUIds(path string) (loadErr error) {
	fd, loadErr := os.OpenFile(path, os.O_RDONLY, 0600)
	defer func() {
		closeErr := fd.Close()
		if closeErr != nil {
			flog.Flogger().Error().Msgf("closing file %s after loading rsa key pairs uuids failed: %s", path, closeErr)
			if loadErr == nil {
				loadErr = closeErr
			} else {
				loadErr = fmt.Errorf("%s\nclosing file %s after loading rsa key pairs uuids failed: %s", loadErr, path, closeErr)
			}
		}
	}()

	if loadErr != nil {
		flog.Flogger().Error().Msgf("opening file %s during loading rsa key pairs uuids failed: %s", path, loadErr)
		loadErr = fmt.Errorf("opening file %s during loading rsa key pairs uuids failed: %s", path, loadErr)
		return loadErr
	}

	reader := bufio.NewReader(fd)

	for {
		line, err := reader.ReadString('\n')
		uuid := strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				if len(uuid) > 0 {
					globalRSAKeyPairsUUIDs = append(globalRSAKeyPairsUUIDs, uuid)
				}
				break
			}

			flog.Flogger().Error().Msgf("reading file %s by line failed: %s", path, err)
			loadErr = fmt.Errorf("reading file %s by line failed: %s", path, err)
			return loadErr
		}
		globalRSAKeyPairsUUIDs = append(globalRSAKeyPairsUUIDs, uuid)
	}

	if len(globalRSAKeyPairsUUIDs) == 0 {
		flog.Flogger().Error().Msgf("the file %s is empty", path)
		loadErr = fmt.Errorf("the file %s is empty", path)
	}

	return loadErr
}
