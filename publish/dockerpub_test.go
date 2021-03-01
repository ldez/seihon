package publish

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDockerPub(t *testing.T) {
	targets, err := GetTargetedArchitectures([]string{"arm.v6", "arm.v7", "arm.v8", "amd64", "386"})
	require.NoError(t, err)

	pub, err := NewDockerPub("image", []string{"v666", "latest"}, "alpine:3.13", targets, "../tmpl.Dockerfile")
	require.NoError(t, err)

	defer func() { _ = pub.Clean(false) }()

	require.Len(t, pub.builds, 5)
	assert.Equal(t, strings.Split("docker build -t image:v666-386 -t image:latest-386 -f linux-386-.Dockerfile .", " "), pub.builds[0].Args)
	assert.Equal(t, strings.Split("docker build -t image:v666-amd64 -t image:latest-amd64 -f linux-amd64-.Dockerfile .", " "), pub.builds[1].Args)
	assert.Equal(t, strings.Split("docker build -t image:v666-arm.v6 -t image:latest-arm.v6 -f linux-arm-6.Dockerfile .", " "), pub.builds[2].Args)
	assert.Equal(t, strings.Split("docker build -t image:v666-arm.v7 -t image:latest-arm.v7 -f linux-arm-7.Dockerfile .", " "), pub.builds[3].Args)
	assert.Equal(t, strings.Split("docker build -t image:v666-arm.v8 -t image:latest-arm.v8 -f linux-arm64-.Dockerfile .", " "), pub.builds[4].Args)

	require.Len(t, pub.push, 10)
	assert.Equal(t, strings.Split("docker push image:v666-386", " "), pub.push[0].Args)
	assert.Equal(t, strings.Split("docker push image:latest-386", " "), pub.push[1].Args)
	assert.Equal(t, strings.Split("docker push image:v666-amd64", " "), pub.push[2].Args)
	assert.Equal(t, strings.Split("docker push image:latest-amd64", " "), pub.push[3].Args)
	assert.Equal(t, strings.Split("docker push image:v666-arm.v6", " "), pub.push[4].Args)
	assert.Equal(t, strings.Split("docker push image:latest-arm.v6", " "), pub.push[5].Args)
	assert.Equal(t, strings.Split("docker push image:v666-arm.v7", " "), pub.push[6].Args)
	assert.Equal(t, strings.Split("docker push image:latest-arm.v7", " "), pub.push[7].Args)
	assert.Equal(t, strings.Split("docker push image:v666-arm.v8", " "), pub.push[8].Args)
	assert.Equal(t, strings.Split("docker push image:latest-arm.v8", " "), pub.push[9].Args)
}

func TestDockerPub_Execute(t *testing.T) {
	targets, err := GetTargetedArchitectures([]string{"arm.v6", "arm.v7", "arm.v8", "amd64", "386"})
	require.NoError(t, err)

	pub, err := NewDockerPub("image", []string{"v666", "latest"}, "alpine:3.13", targets, "../tmpl.Dockerfile")
	require.NoError(t, err)

	defer func() { _ = pub.Clean(false) }()

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

	assert.Equal(t, `docker build -t image:v666-386 -t image:latest-386 -f linux-386-.Dockerfile .
docker build -t image:v666-amd64 -t image:latest-amd64 -f linux-amd64-.Dockerfile .
docker build -t image:v666-arm.v6 -t image:latest-arm.v6 -f linux-arm-6.Dockerfile .
docker build -t image:v666-arm.v7 -t image:latest-arm.v7 -f linux-arm-7.Dockerfile .
docker build -t image:v666-arm.v8 -t image:latest-arm.v8 -f linux-arm64-.Dockerfile .
docker push image:v666-386
docker push image:latest-386
docker push image:v666-amd64
docker push image:latest-amd64
docker push image:v666-arm.v6
docker push image:latest-arm.v6
docker push image:v666-arm.v7
docker push image:latest-arm.v7
docker push image:v666-arm.v8
docker push image:latest-arm.v8
`, string(out))
}
