package util

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/moby/client"
	mobyClient "github.com/moby/moby/client"
)

func PullImageIfNotExists(client *client.Client, ctx context.Context, image string) error {
	output, err := client.ImagePull(ctx, image, mobyClient.ImagePullOptions{})
	if err != nil {
		return err
	}

	defer func(output mobyClient.ImagePullResponse) {
		err := output.Close()
		if err != nil {
			fmt.Printf("Can't close output.")
		}
	}(output)

	type progressDetail struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	}

	type progressMsg struct {
		Status         string         `json:"status"`
		ID             string         `json:"id"`
		ProgressDetail progressDetail `json:"progressDetail"`
	}

	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		var message progressMsg

		line := scanner.Bytes()
		if err := json.Unmarshal(line, &message); err != nil {
			continue
		}

		if message.ProgressDetail.Total > 0 {
			percent := float64(message.ProgressDetail.Current) / float64(message.ProgressDetail.Total) * 100
			barWidth := 20
			filled := int(percent / 100 * float64(barWidth))

			fmt.Printf("\r[%s%s] %.1f%% %s",
				repeatRune('█', filled),          // barra preenchida
				repeatRune('░', barWidth-filled), // barra “vazia”
				percent,
				message.ID,
			)
		} else {
			fmt.Printf("\r%s: %s", message.ID, message.Status)
		}
	}
	fmt.Println()
	return nil
}

func repeatRune(r rune, count int) string {
	s := make([]rune, count)
	for i := range s {
		s[i] = r
	}
	return string(s)
}
