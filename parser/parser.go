package parser

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/emicklei/proto"
	"github.com/pkg/errors"
	"github.com/serenize/snaker"
)

// Parse parses the proto file.
func Parse(pf *proto.Proto) (*File, error) {
	f := &File{}
	var packageName string
	var messages []*proto.Message
	for _, element := range pf.Elements {
		switch ele := element.(type) {
		case *proto.Package:
			packageName = ele.Name
		case *proto.Message:
			messages = append(messages, ele)
		}
	}

	if packageName == "" {
		return nil, errors.New("package name not found in proto file")
	}

	f.PackageName = packageName
	mapTypes := make(map[string]*Map)
	messageTypes := make(map[string]*Message)
	genMapName := func(mele *proto.MapField) string {
		var outp bytes.Buffer
		outp.WriteString("I")
		outp.WriteString(snaker.SnakeToCamel(mele.KeyType))
		outp.WriteString(snaker.SnakeToCamel(mele.Type))
		outp.WriteString("Map")
		return outp.String()
	}

	for _, message := range messages {
		var msg Message
		msg.Name = message.Name
		msg.InterName = "I" + message.Name
		if message.Comment != nil {
			msg.Comment = strings.TrimSpace(message.Comment.Message())
		}

		messageTypes[message.Name] = &msg
		for _, melement := range message.Elements {
			switch mele := melement.(type) {
			case *proto.NormalField:
				var comment string
				if mele.Comment != nil {
					comment = strings.TrimSpace(mele.Comment.Message())
				}

				msg.Fields = append(msg.Fields, Field{
					Name:      mele.Name,
					CamelName: snaker.SnakeToCamel(mele.Name),
					Comment:   comment,
					Type:      mele.Type,
				})
			case *proto.MapField:
				var comment string
				if mele.Comment != nil {
					comment = strings.TrimSpace(mele.Comment.Message())
				}

				mapName := genMapName(mele)
				mt, ok := mapTypes[mapName]
				if !ok {
					f.Maps = append(f.Maps, Map{
						Key:      mele.KeyType,
						Value:    mele.Type,
						TypeName: mapName,
						ValuePtr: mele.Type,
					})
					mt = &f.Maps[len(f.Maps)-1]
				}

				msg.Fields = append(msg.Fields, Field{
					Name:      mele.Name,
					CamelName: snaker.SnakeToCamel(mele.Name),
					Comment:   comment,
					Type:      mele.Type,
					Map:       mt,
				})
			}
		}

		f.Messages = append(f.Messages, msg)
	}

	sort.Slice(f.Messages, func(i int, j int) bool {
		return strings.Compare(f.Messages[i].Name, f.Messages[j].Name) == -1
	})

	sort.Slice(f.Maps, func(i int, j int) bool {
		return strings.Compare(f.Maps[i].TypeName, f.Maps[j].TypeName) == -1
	})

	for k, ma := range f.Maps {
		fmt.Println(ma.Value)
		it, recog := messageTypes[ma.Value]
		fmt.Println(recog)
		if recog {
			fmt.Println(it.InterName)
			f.Maps[k].ValuePtr = "*" + ma.Value
			f.Maps[k].Value = it.InterName
		}
	}

	return f, nil
}
