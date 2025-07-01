package internal

import (
	"context"
	"io"
	"os"

	proto "github.com/codek7-services/codek7-tui/pkg/pb"
)

func UploadVideo(client proto.RepoServiceClient, filePath, title, description, userID string) error {
	stream, err := client.UploadVideo(context.TODO())
	if err != nil {
		return err
	}

	// Send metadata
	err = stream.Send(&proto.UploadVideoRequest{
		Data: &proto.UploadVideoRequest_Metadata{
			Metadata: &proto.VideoMetadata{
				UserId:      userID,
				Title:       title,
				Description: description,
				FileName:    filePath,
				FileSize:    0, // optional
			},
		},
	})
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 1024*64)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		err = stream.Send(&proto.UploadVideoRequest{
			Data: &proto.UploadVideoRequest_Chunk{
				Chunk: &proto.VideoChunk{
					Data:        buf[:n],
					ChunkNumber: 1,
				},
			},
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	return err
}
