package listener

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"
	"time"
	"github.com/spf13/viper"
"github.com/sayden/gubsub/config"
)

func NewFileListener(filePath types.FileListener, topic *string) (*uuid.UUID, error) {
	f, err := os.Create(filePath.Path)
	if err != nil {
		return nil, err
	}

	//Creates new listener
	l := types.NewListener(topic)
	dispatcher.AddListenerToTopic(l)

	go launchFileWriterGoroutine(l, f, &filePath)

	return &l.ID, nil
}

func launchFileWriterGoroutine(l *types.Listener, f *os.File, filePath *types.FileListener) {

	var sync, quit chan bool
	go SyncWrite(sync, quit)

	for {
		select {
		case m := <-l.Ch:
			_, err := f.Write(append(*m.Data, []byte("\n")...))
			if err != nil {
				log.WithFields(log.Fields{
					"file":  filePath.Path,
					"error": err.Error(),
				}).Error("Could not write to file")
			}
		case <-sync:
			f.Sync()
		}
	}

	quit <- true
	close(sync)

	f.Sync()
	f.Close()
}

func SyncWrite(c chan bool, q chan bool) {
	for {
		select {
		case <-q:
			return
		default:
			time.Sleep(time.Duration(viper.GetInt(config.WRITE_TO_FILE_DELAY)) * time.Second)
			c <- true
		}
	}
}
