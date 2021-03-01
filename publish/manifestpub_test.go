package publish

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManifestPub(t *testing.T) {
	pub, err := NewManifestPub("image", "v666", availableArchitectures)
	require.NoError(t, err)

	require.Len(t, pub.manifestAnnotate, 6)

	assert.Equal(t, strings.Split("docker manifest annotate image:v666 image:v666-386 --os=linux --arch=386", " "), pub.manifestAnnotate[0].Args)
	assert.Equal(t, strings.Split("docker manifest annotate image:v666 image:v666-amd64 --os=linux --arch=amd64", " "), pub.manifestAnnotate[1].Args)
	assert.Equal(t, strings.Split("docker manifest annotate image:v666 image:v666-arm.v5 --os=linux --arch=arm --variant=v5", " "), pub.manifestAnnotate[2].Args)
	assert.Equal(t, strings.Split("docker manifest annotate image:v666 image:v666-arm.v6 --os=linux --arch=arm --variant=v6", " "), pub.manifestAnnotate[3].Args)
	assert.Equal(t, strings.Split("docker manifest annotate image:v666 image:v666-arm.v7 --os=linux --arch=arm --variant=v7", " "), pub.manifestAnnotate[4].Args)
	assert.Equal(t, strings.Split("docker manifest annotate image:v666 image:v666-arm.v8 --os=linux --arch=arm64 --variant=v8", " "), pub.manifestAnnotate[5].Args)

	assert.Equal(t, []string{
		"docker", "manifest", "create", "--amend",
		"image:v666",
		"image:v666-386",
		"image:v666-amd64",
		"image:v666-arm.v5",
		"image:v666-arm.v6",
		"image:v666-arm.v7",
		"image:v666-arm.v8",
	}, pub.manifestCreate.Args)

	assert.Equal(t, strings.Split("docker manifest push --purge image:v666", " "), pub.manifestPush.Args)
}

func TestManifestPub_Execute(t *testing.T) {
	pub, err := NewManifestPub("image", "v666", availableArchitectures)
	require.NoError(t, err)

	// catch stdout
	backupStdout := os.Stdout
	defer func() { os.Stdout = backupStdout }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = pub.Execute(true)
	require.NoError(t, err)

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = backupStdout

	assert.Equal(t, `docker manifest create --amend image:v666 image:v666-386 image:v666-amd64 image:v666-arm.v5 image:v666-arm.v6 image:v666-arm.v7 image:v666-arm.v8
docker manifest annotate image:v666 image:v666-386 --os=linux --arch=386
docker manifest annotate image:v666 image:v666-amd64 --os=linux --arch=amd64
docker manifest annotate image:v666 image:v666-arm.v5 --os=linux --arch=arm --variant=v5
docker manifest annotate image:v666 image:v666-arm.v6 --os=linux --arch=arm --variant=v6
docker manifest annotate image:v666 image:v666-arm.v7 --os=linux --arch=arm --variant=v7
docker manifest annotate image:v666 image:v666-arm.v8 --os=linux --arch=arm64 --variant=v8
docker manifest push --purge image:v666
`, string(out))
}
