package profiles

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path"

	"aeroheart.io/hyprdisp/hyprpanel"
	"aeroheart.io/hyprdisp/sys"
	"github.com/pelletier/go-toml/v2"
)

func (s defaultService) applyPanels(ctx context.Context) error {
	var (
		logger   *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		filePath string      = path.Join(".", "var", "panels.toml")
		data     []byte
		err      error
	)

	data, err = os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var profile panelProfile
	err = toml.Unmarshal(data, &profile)
	if err != nil {
		return err
	}

	logger.Printf("Panel config: %+v", profile)

	var layout hyprpanel.BarLayout = make(hyprpanel.BarLayout)

	layout.Set("0", hyprpanel.BarWidgetConfig{
		L: profile["main"].L,
		R: profile["main"].R,
		M: profile["main"].M,
	})
	layout.Set("1", hyprpanel.BarWidgetConfig{
		L: profile["sub"].L,
		R: profile["sub"].R,
		M: profile["sub"].M,
	})
	layout.Set("2", hyprpanel.BarWidgetConfig{
		L: profile["sub"].L,
		R: profile["sub"].R,
		M: profile["sub"].M,
	})
	layout.Set("3", hyprpanel.BarWidgetConfig{
		L: profile["sub"].L,
		R: profile["sub"].R,
		M: profile["sub"].M,
	})

	var jsondata []byte
	jsondata, err = json.Marshal(layout)
	if err != nil {
		logger.Printf("OH SHIT: %+v", err)
	}
	logger.Printf("Panel config asJSON: %+v", string(jsondata))
	return nil
}
