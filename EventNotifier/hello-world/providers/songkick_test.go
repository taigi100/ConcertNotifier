package providers

import (
	"reflect"
	"testing"
)

func TestPovider(t *testing.T) {
	t.Run("Unable to construct provider", func(t *testing.T) {
		_, err := NewSongKickConnector()

		if err != nil {
			t.Fatal("Error constructing SongKick provider object")
		}
	})

	t.Run("Wrong ConfigObjectType", func(t *testing.T) {
		provider, _ := NewSongKickConnector()

		_, err := provider.GetAllEvents("random object")

		if err == nil {
			t.Fatalf("Wrong config object should return error, got %v", err)
		}
	})

	t.Run("Should get some artists for myself", func(t *testing.T) {
		provider, _ := NewSongKickConnector()
		ret, err := provider.getAllFollowedArtists(SongKickConfig{Username: "robert-cazacu"})

		if err != nil {
			t.Error(err)
		}

		if len(ret) == 0 {
			t.Fatal("Artists array is empty")
		}
	})

	t.Run("Should fail for unknown user", func(t *testing.T) {
		provider, _ := NewSongKickConnector()
		_, err := provider.getAllFollowedArtists(SongKickConfig{Username: "{uersras}"})

		if err == nil {
			t.Fatal("getAllArtistsFollowed did not throw err for unknown user")
		}
	})

}

func Test_songKickConnector_getAllArtistsConcerts(t *testing.T) {
	type args struct {
		artist songKickArtist
	}
	// connector, _ := NewSongKickConnector()
	tests := []struct {
		name     string
		provider songKickConnector
		args     args
		want     []Event
		wantErr  bool
	}{
		// TODO: Add test cases.
		// ! How do I test this considering I'm getting a huge response
		// ! and the reponse is dynamic :-?
		// {
		// 	name:     "Get concerts for inexisting artist",
		// 	provider: connector,
		// 	args:     args{songKickArtist{id: "277064"}},
		// 	want:     nil,
		// 	wantErr:  false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.provider.getAllArtistsConcerts(tt.args.artist)
			if (err != nil) != tt.wantErr {
				t.Errorf("songKickConnector.getAllArtistsConcerts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("songKickConnector.getAllArtistsConcerts() = %v, want %v", got, tt.want)
			}
		})
	}
}
