package musicbrainz

import (
	"context"
	"fmt"
)

type MBService struct {
	repo Storage
}

func NewMBService(repo Storage) *MBService {
	return &MBService{repo: repo}
}

type TrackInfo struct {
	RecordingMBID string
	RecordingName string
	ArtistName    string
	ArtistMBID    string
	Genres        []Genre
}

func (s TrackInfo) IsEmpty() bool {
	return s.RecordingMBID == "" && len(s.Genres) == 0
}

func (mb *MBService) GetMBIDFromArtistAndTrackName(context context.Context, artistName string, trackName string) (TrackInfo, error) {
	topCandidate, err := mb.repo.Recording.GetClosestRecording(context, trackName, artistName)
	if err != nil {
		return TrackInfo{}, fmt.Errorf("GetMBIDFromArtistAndTrackName error: %v", err)
	}
	genresMap, err := mb.repo.Genre.GetGenreByRecordings(context, []string{topCandidate.RecordingMBID})
	if err != nil {
		return TrackInfo{}, fmt.Errorf("GetMBIDFromArtistAndTrackName error: %v", err)
	}
	genres := genresMap[topCandidate.RecordingMBID]
	return TrackInfo{
		RecordingMBID: topCandidate.RecordingMBID,
		RecordingName: topCandidate.RecordingName,
		ArtistName:    topCandidate.ArtistName,
		ArtistMBID:    topCandidate.ArtistMBID,
		Genres:        genres,
	}, nil

}

func (mb *MBService) GetMBIDsFromArtistAndTrackNames(context context.Context, artistNames []string, trackNames []string) ([]TrackInfo, error) {
	candidates, err := mb.repo.Recording.GetClosestRecordings(context, trackNames, artistNames)
	if err != nil {
		return nil, fmt.Errorf("GetMBIDsFromArtistAndTrackNames error: %v", err)
	}

	recordingMBIDs := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		recordingMBIDs = append(recordingMBIDs, candidate.RecordingMBID)
	}

	trackInfos := make([]TrackInfo, len(candidates))
	for i := range len(candidates) {
		trackInfos[i] = TrackInfo{
			RecordingMBID: candidates[i].RecordingMBID,
			RecordingName: candidates[i].RecordingName,
			ArtistName:    candidates[i].ArtistName,
			ArtistMBID:    candidates[i].ArtistMBID,
		}
	}

	return trackInfos, nil
}

func (mb *MBService) GetGenresByRecordingMBIDs(context context.Context, recordingMBIDs []string) (map[string][]Genre, error) {
	return mb.repo.Genre.GetGenreByRecordings(context, recordingMBIDs)
}

func (mb *MBService) GetGenresByRecordingMBID(context context.Context, recordingMBID string) ([]Genre, error) {
	genresMap, err := mb.repo.Genre.GetGenreByRecordings(context, []string{recordingMBID})
	if err != nil {
		return nil, err
	}
	return genresMap[recordingMBID], nil
}

func (mb *MBService) IsValidMBIDForRecord(context context.Context, mbid string) (bool, error) {
	return mb.repo.Recording.DoesRecordExistByMBID(context, mbid)
}
