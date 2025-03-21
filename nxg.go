package nxg

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// NXGHeader represents the decoded components of an NXG header.
type DecodedNxgHeader struct {
	RandomString   string
	TotalDataParts int64
	TotalParParts  int64
}

// Error messages
var (
	ErrFailedDecodeBase64Header       = fmt.Errorf("failed to decode base64 header")
	ErrUnexpectedHeaderFmt            = fmt.Errorf("unexpected header format")
	ErrFailedToParseTotalDataArticles = fmt.Errorf("failed to parse total data articles")
	ErrParseTotalPar2Articles         = fmt.Errorf("failed to parse total par2 articles")
	ErrUnexpectedHashLength           = fmt.Errorf("unexpected hash length")
)

type PartType string

const (
	PartTypeData PartType = "data"
	PartTypePar2 PartType = "par2"
)

const (
	expectedNxgHeaderHashLength = 64
)

type Header [28]byte

func (h Header) String() string {
	return string(h[:])
}

// GenerateSegmentID generates the segment ID based on the NXG header, article type, and article number.
func (h Header) GenerateSegmentID(partType PartType, partNumber int64) (string, error) {
	// Get the segment hash
	segmentHash, err := h.getSegmentHash(partType, partNumber)
	if err != nil {
		return "", err
	}

	// Format the message ID
	messageID := fmt.Sprintf("%s@%s.%s", segmentHash[:40], segmentHash[40:61], segmentHash[61:64])

	return messageID, nil
}

func (h Header) GetXNxgHeader(
	fileNumber,
	totalFiles int64,
	filename string,
	partType PartType,
	totalDownloadSize int64,
) (string, error) {
	return encrypt(fmt.Sprintf("%v:%v:%v:%v:%v",
		fileNumber,
		totalFiles,
		filename,
		partType,
		totalDownloadSize,
	), h.String())
}

func (h Header) GetObfuscatedSubject(
	partType PartType,
	partNumber int64,
) (string, error) {
	segmentHash, err := h.getSegmentHash(partType, partNumber)
	if err != nil {
		return "", err
	}

	return getSHA256Hash(segmentHash), nil
}

func (h Header) GetObfuscatedPoster(
	partType PartType,
	partNumber int64,
) (string, error) {
	obfuscatedSubject, err := h.GetObfuscatedSubject(partType, partNumber)
	if err != nil {
		return "", err
	}

	posterHash := getSHA256Hash(obfuscatedSubject)
	obfuscatedPoster := posterHash[10:15] + " <" + posterHash[10:25] + "@" + posterHash[30:45] + "." + posterHash[50:53] + ">"

	return obfuscatedPoster, nil
}

func (h Header) getSegmentHash(
	partType PartType,
	partNumber int64,
) (string, error) {
	// Construct the string for hashing
	text := fmt.Sprintf("%v:%v:%v", h.String(), partType, partNumber)

	// Calculate the SHA256 hash
	hasher := sha256.New()
	hasher.Write([]byte(text))
	hashStr := hex.EncodeToString(hasher.Sum(nil))

	// Ensure the hash string length is as expected
	if len(hashStr) != expectedNxgHeaderHashLength {
		return "", fmt.Errorf("%w: %d", ErrUnexpectedHashLength, len(hashStr))
	}

	return hashStr, nil
}

func (h Header) GetNxgLink(
	queryStr map[string][]string,
) string {
	return fmt.Sprintf("nxglnk://?h=%s&%s", h.String(), url.Values(queryStr).Encode())
}

// GenerateNXGHeader generates the NXG header string from its components.
func GenerateNXGHeader(totalDataParts, totalParParts int64) Header {
	shortHeader := randomString(21, 0, true)
	partBytes := []byte(fmt.Sprintf(":%d:%d", totalDataParts, totalParParts))
	fullHeaderBytes := []byte(shortHeader)

	copy(fullHeaderBytes[len(fullHeaderBytes)-len(partBytes):], partBytes)

	var header Header

	base64.StdEncoding.Encode(header[:], fullHeaderBytes)

	return header
}

// DecodeNXGHeader decodes a base64-encoded NXG header into its components.
func DecodeNXGHeader(encodedHeader string) (*DecodedNxgHeader, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedHeader)
	if err != nil {
		return nil, errors.Join(ErrFailedDecodeBase64Header, err)
	}

	parts := strings.Split(string(decodedBytes), ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedHeaderFmt, string(decodedBytes))
	}

	totalDataParts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, errors.Join(ErrFailedToParseTotalDataArticles, err)
	}

	totalParParts, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, errors.Join(ErrParseTotalPar2Articles, err)
	}

	header := &DecodedNxgHeader{
		RandomString:   parts[0],
		TotalDataParts: totalDataParts,
		TotalParParts:  totalParParts,
	}

	return header, nil
}
