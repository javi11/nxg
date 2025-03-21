package nxg

import (
	"errors"
	"testing"
)

func TestDecodeNXGHeader(t *testing.T) {
	tests := []struct {
		wantHeader    *DecodedNxgHeader
		wantErr       error
		name          string
		encodedHeader string
	}{
		{
			name:          "valid header",
			encodedHeader: "aGVsbG86NToxMA==",
			wantHeader: &DecodedNxgHeader{
				RandomString:   "hello",
				TotalDataParts: 5,
				TotalParParts:  10,
			},
			wantErr: nil,
		},
		{
			name:          "invalid base64 encoding",
			encodedHeader: "SGVsbG8=:5:10!",
			wantHeader:    nil,
			wantErr:       ErrFailedDecodeBase64Header,
		},
		{
			name:          "unexpected header format",
			encodedHeader: "aGVsbG86NQ==",
			wantHeader:    nil,
			wantErr:       ErrUnexpectedHeaderFmt,
		},
		{
			name:          "failed to parse total data articles",
			encodedHeader: "aGVsbG86Zml2ZToxMA==",
			wantHeader:    nil,
			wantErr:       ErrFailedToParseTotalDataArticles,
		},
		{
			name:          "failed to parse total par2 articles",
			encodedHeader: "aGVsbG86NTp0ZW4=",
			wantHeader:    nil,
			wantErr:       ErrParseTotalPar2Articles,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHeader, gotErr := DecodeNXGHeader(tt.encodedHeader)
			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Errorf("DecodeNXGHeader() error = %v, wantErr %v", gotErr, tt.wantErr)
				}

				return
			}

			if gotErr != nil {
				t.Errorf("DecodeNXGHeader() unexpected error = %v", gotErr)
			}

			if !headerEqual(gotHeader, tt.wantHeader) {
				t.Errorf("DecodeNXGHeader() = %v, want %v", gotHeader, tt.wantHeader)
			}
		})
	}
}

func TestGenerateNXGHeader(t *testing.T) {
	totalDataParts := int64(5)
	totalParParts := int64(10)

	header := GenerateNXGHeader(totalDataParts, totalParParts)

	decodedHeader, err := DecodeNXGHeader(header.String())
	if err != nil {
		t.Errorf("DecodeNXGHeader() unexpected error = %v", err)
	}

	if decodedHeader.TotalDataParts != totalDataParts {
		t.Errorf("TotalDataParts = %v, want %v", decodedHeader.TotalDataParts, totalDataParts)
	}

	if decodedHeader.TotalParParts != totalParParts {
		t.Errorf("TotalParParts = %v, want %v", decodedHeader.TotalParParts, totalParParts)
	}
}

func TestNxgHeader_GenerateSegmentID(t *testing.T) {
	nxgHeader := GenerateNXGHeader(5, 10)
	partType := PartTypeData
	partNumber := int64(3)

	segmentID, err := nxgHeader.GenerateSegmentID(partType, partNumber)
	if err != nil {
		t.Errorf("GenerateSegmentID() unexpected error = %v", err)
	}

	if segmentID == "" {
		t.Error("GenerateSegmentID() returned empty string")
	}
}

func TestNxgHeader_GetXNxgHeader(t *testing.T) {
	nxgHeader := GenerateNXGHeader(5, 10)
	fileNumber := int64(1)
	totalFiles := int64(3)
	filename := "test.txt"
	partType := PartTypeData
	totalDownloadSize := int64(1024)

	xNxgHeader, err := nxgHeader.GetXNxgHeader(fileNumber, totalFiles, filename, partType, totalDownloadSize)
	if err != nil {
		t.Errorf("GetXNxgHeader() unexpected error = %v", err)
	}

	if xNxgHeader == "" {
		t.Error("GetXNxgHeader() returned empty string")
	}
}

func TestNxgHeader_GetObfuscatedSubject(t *testing.T) {
	nxgHeader := GenerateNXGHeader(5, 10)
	partType := PartTypeData
	partNumber := int64(3)

	obfuscatedSubject, err := nxgHeader.GetObfuscatedSubject(partType, partNumber)
	if err != nil {
		t.Errorf("GetObfuscatedSubject() unexpected error = %v", err)
	}

	if obfuscatedSubject == "" {
		t.Error("GetObfuscatedSubject() returned empty string")
	}
}

func TestNxgHeader_GetObfuscatedPoster(t *testing.T) {
	nxgHeader := GenerateNXGHeader(5, 10)
	partType := PartTypeData
	partNumber := int64(3)

	obfuscatedPoster, err := nxgHeader.GetObfuscatedPoster(partType, partNumber)
	if err != nil {
		t.Errorf("GetObfuscatedPoster() unexpected error = %v", err)
	}

	if obfuscatedPoster == "" {
		t.Error("GetObfuscatedPoster() returned empty string")
	}
}

func headerEqual(h1, h2 *DecodedNxgHeader) bool {
	if h1 == nil && h2 == nil {
		return true
	}

	if h1 == nil || h2 == nil {
		return false
	}

	return h1.RandomString == h2.RandomString &&
		h1.TotalDataParts == h2.TotalDataParts &&
		h1.TotalParParts == h2.TotalParParts
}
