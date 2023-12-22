/*
 * Copyright 2023 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cast_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/lvan100/cast"
	"github.com/lvan100/cast/internal/assert"
)

type TwitterStruct struct {
	Statuses []struct {
		Coordinates interface{} `json:"coordinates"`
		Favorited   bool        `json:"favorited"`
		Truncated   bool        `json:"truncated"`
		CreatedAt   string      `json:"created_at"`
		IdStr       string      `json:"id_str"`
		Entities    struct {
			Urls     []interface{} `json:"urls"`
			Hashtags []struct {
				Text    string    `json:"text"`
				Indices []float64 `json:"indices"`
			} `json:"hashtags"`
			UserMentions []interface{} `json:"user_mentions"`
		} `json:"entities"`
		InReplyToUserIdStr interface{} `json:"in_reply_to_user_id_str"`
		Contributors       interface{} `json:"contributors"`
		Text               string      `json:"text"`
		Metadata           struct {
			IsoLanguageCode string `json:"iso_language_code"`
			ResultType      string `json:"result_type"`
		} `json:"metadata"`
		RetweetCount         float64     `json:"retweet_count"`
		InReplyToStatusIdStr interface{} `json:"in_reply_to_status_id_str"`
		Id                   float64     `json:"id"`
		Geo                  interface{} `json:"geo"`
		Retweeted            bool        `json:"retweeted"`
		InReplyToUserId      interface{} `json:"in_reply_to_user_id"`
		Place                interface{} `json:"place"`
		User                 struct {
			ProfileSidebarFillColor   string      `json:"profile_sidebar_fill_color"`
			ProfileSidebarBorderColor string      `json:"profile_sidebar_border_color"`
			ProfileBackgroundTile     bool        `json:"profile_background_tile"`
			Name                      string      `json:"name"`
			ProfileImageUrl           string      `json:"profile_image_url"`
			CreatedAt                 string      `json:"created_at"`
			Location                  string      `json:"location"`
			FollowRequestSent         interface{} `json:"follow_request_sent"`
			ProfileLinkColor          string      `json:"profile_link_color"`
			IsTranslator              bool        `json:"is_translator"`
			IdStr                     string      `json:"id_str"`
			Entities                  struct {
				Url struct {
					Urls []struct {
						ExpandedUrl interface{} `json:"expanded_url"`
						Url         *string     `json:"url"`
						Indices     []float64   `json:"indices"`
					} `json:"urls"`
				} `json:"url"`
				Description struct {
					Urls []interface{} `json:"urls"`
				} `json:"description"`
			} `json:"entities"`
			DefaultProfile                 bool        `json:"default_profile"`
			ContributorsEnabled            bool        `json:"contributors_enabled"`
			FavouritesCount                float64     `json:"favourites_count"`
			Url                            *string     `json:"url"`
			ProfileImageUrlHttps           string      `json:"profile_image_url_https"`
			UtcOffset                      float64     `json:"utc_offset"`
			Id                             float64     `json:"id"`
			ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
			ListedCount                    float64     `json:"listed_count"`
			ProfileTextColor               string      `json:"profile_text_color"`
			Lang                           string      `json:"lang"`
			FollowersCount                 float64     `json:"followers_count"`
			Protected                      bool        `json:"protected"`
			Notifications                  interface{} `json:"notifications"`
			ProfileBackgroundImageUrlHttps string      `json:"profile_background_image_url_https"`
			ProfileBackgroundColor         string      `json:"profile_background_color"`
			Verified                       bool        `json:"verified"`
			GeoEnabled                     bool        `json:"geo_enabled"`
			TimeZone                       string      `json:"time_zone"`
			Description                    string      `json:"description"`
			DefaultProfileImage            bool        `json:"default_profile_image"`
			ProfileBackgroundImageUrl      string      `json:"profile_background_image_url"`
			StatusesCount                  float64     `json:"statuses_count"`
			FriendsCount                   float64     `json:"friends_count"`
			Following                      interface{} `json:"following"`
			ShowAllInlineMedia             bool        `json:"show_all_inline_media"`
			ScreenName                     string      `json:"screen_name"`
		} `json:"user"`
		InReplyToScreenName interface{} `json:"in_reply_to_screen_name"`
		Source              string      `json:"source"`
		InReplyToStatusId   interface{} `json:"in_reply_to_status_id"`
	} `json:"statuses"`
	SearchMetadata struct {
		MaxId       float64 `json:"max_id"`
		SinceId     float64 `json:"since_id"`
		RefreshUrl  string  `json:"refresh_url"`
		NextResults string  `json:"next_results"`
		Count       float64 `json:"count"`
		CompletedIn float64 `json:"completed_in"`
		SinceIdStr  string  `json:"since_id_str"`
		Query       string  `json:"query"`
		MaxIdStr    string  `json:"max_id_str"`
	} `json:"search_metadata"`
}

const TwitterJson = `{
  "statuses": [
    {
      "coordinates": null,
      "favorited": false,
      "truncated": false,
      "created_at": "Mon Sep 24 03:35:21 +0000 2012",
      "id_str": "250075927172759552",
      "entities": {
        "urls": [
 
        ],
        "hashtags": [
          {
            "text": "freebandnames",
            "indices": [
              20,
              34
            ]
          }
        ],
        "user_mentions": [
 
        ]
      },
      "in_reply_to_user_id_str": null,
      "contributors": null,
      "text": "Aggressive Ponytail #freebandnames",
      "metadata": {
        "iso_language_code": "en",
        "result_type": "recent"
      },
      "retweet_count": 0,
      "in_reply_to_status_id_str": null,
      "id": 250075927172759552,
      "geo": null,
      "retweeted": false,
      "in_reply_to_user_id": null,
      "place": null,
      "user": {
        "profile_sidebar_fill_color": "DDEEF6",
        "profile_sidebar_border_color": "C0DEED",
        "profile_background_tile": false,
        "name": "Sean Cummings",
        "profile_image_url": "https://a0.twimg.com/profile_images/2359746665/1v6zfgqo8g0d3mk7ii5s_normal.jpeg",
        "created_at": "Mon Apr 26 06:01:55 +0000 2010",
        "location": "LA, CA",
        "follow_request_sent": null,
        "profile_link_color": "0084B4",
        "is_translator": false,
        "id_str": "137238150",
        "entities": {
          "url": {
            "urls": [
              {
                "expanded_url": null,
                "url": "",
                "indices": [
                  0,
                  0
                ]
              }
            ]
          },
          "description": {
            "urls": [
 
            ]
          }
        },
        "default_profile": true,
        "contributors_enabled": false,
        "favourites_count": 0,
        "url": null,
        "profile_image_url_https": "https://si0.twimg.com/profile_images/2359746665/1v6zfgqo8g0d3mk7ii5s_normal.jpeg",
        "utc_offset": -28800,
        "id": 137238150,
        "profile_use_background_image": true,
        "listed_count": 2,
        "profile_text_color": "333333",
        "lang": "en",
        "followers_count": 70,
        "protected": false,
        "notifications": null,
        "profile_background_image_url_https": "https://si0.twimg.com/images/themes/theme1/bg.png",
        "profile_background_color": "C0DEED",
        "verified": false,
        "geo_enabled": true,
        "time_zone": "Pacific Time (US & Canada)",
        "description": "Born 330 Live 310",
        "default_profile_image": false,
        "profile_background_image_url": "https://a0.twimg.com/images/themes/theme1/bg.png",
        "statuses_count": 579,
        "friends_count": 110,
        "following": null,
        "show_all_inline_media": false,
        "screen_name": "sean_cummings"
      },
      "in_reply_to_screen_name": null,
      "source": "<a href=\"//itunes.apple.com/us/app/twitter/id409789998?mt=12%5C%22\" rel=\"\\\"nofollow\\\"\">Twitter for Mac</a>",
      "in_reply_to_status_id": null
    },
    {
      "coordinates": null,
      "favorited": false,
      "truncated": false,
      "created_at": "Fri Sep 21 23:40:54 +0000 2012",
      "id_str": "249292149810667520",
      "entities": {
        "urls": [
 
        ],
        "hashtags": [
          {
            "text": "FreeBandNames",
            "indices": [
              20,
              34
            ]
          }
        ],
        "user_mentions": [
 
        ]
      },
      "in_reply_to_user_id_str": null,
      "contributors": null,
      "text": "Thee Namaste Nerdz. #FreeBandNames",
      "metadata": {
        "iso_language_code": "pl",
        "result_type": "recent"
      },
      "retweet_count": 0,
      "in_reply_to_status_id_str": null,
      "id": 249292149810667520,
      "geo": null,
      "retweeted": false,
      "in_reply_to_user_id": null,
      "place": null,
      "user": {
        "profile_sidebar_fill_color": "DDFFCC",
        "profile_sidebar_border_color": "BDDCAD",
        "profile_background_tile": true,
        "name": "Chaz Martenstein",
        "profile_image_url": "https://a0.twimg.com/profile_images/447958234/Lichtenstein_normal.jpg",
        "created_at": "Tue Apr 07 19:05:07 +0000 2009",
        "location": "Durham, NC",
        "follow_request_sent": null,
        "profile_link_color": "0084B4",
        "is_translator": false,
        "id_str": "29516238",
        "entities": {
          "url": {
            "urls": [
              {
                "expanded_url": null,
                "url": "https://bullcityrecords.com/wnng/",
                "indices": [
                  0,
                  32
                ]
              }
            ]
          },
          "description": {
            "urls": [
 
            ]
          }
        },
        "default_profile": false,
        "contributors_enabled": false,
        "favourites_count": 8,
        "url": "https://bullcityrecords.com/wnng/",
        "profile_image_url_https": "https://si0.twimg.com/profile_images/447958234/Lichtenstein_normal.jpg",
        "utc_offset": -18000,
        "id": 29516238,
        "profile_use_background_image": true,
        "listed_count": 118,
        "profile_text_color": "333333",
        "lang": "en",
        "followers_count": 2052,
        "protected": false,
        "notifications": null,
        "profile_background_image_url_https": "https://si0.twimg.com/profile_background_images/9423277/background_tile.bmp",
        "profile_background_color": "9AE4E8",
        "verified": false,
        "geo_enabled": false,
        "time_zone": "Eastern Time (US & Canada)",
        "description": "You will come to Durham, North Carolina. I will sell you some records then, here in Durham, North Carolina. Fun will happen.",
        "default_profile_image": false,
        "profile_background_image_url": "https://a0.twimg.com/profile_background_images/9423277/background_tile.bmp",
        "statuses_count": 7579,
        "friends_count": 348,
        "following": null,
        "show_all_inline_media": true,
        "screen_name": "bullcityrecords"
      },
      "in_reply_to_screen_name": null,
      "source": "web",
      "in_reply_to_status_id": null
    },
    {
      "coordinates": null,
      "favorited": false,
      "truncated": false,
      "created_at": "Fri Sep 21 23:30:20 +0000 2012",
      "id_str": "249289491129438208",
      "entities": {
        "urls": [
 
        ],
        "hashtags": [
          {
            "text": "freebandnames",
            "indices": [
              29,
              43
            ]
          }
        ],
        "user_mentions": [
 
        ]
      },
      "in_reply_to_user_id_str": null,
      "contributors": null,
      "text": "Mexican Heaven, Mexican Hell #freebandnames",
      "metadata": {
        "iso_language_code": "en",
        "result_type": "recent"
      },
      "retweet_count": 0,
      "in_reply_to_status_id_str": null,
      "id": 249289491129438208,
      "geo": null,
      "retweeted": false,
      "in_reply_to_user_id": null,
      "place": null,
      "user": {
        "profile_sidebar_fill_color": "99CC33",
        "profile_sidebar_border_color": "829D5E",
        "profile_background_tile": false,
        "name": "Thomas John Wakeman",
        "profile_image_url": "https://a0.twimg.com/profile_images/2219333930/Froggystyle_normal.png",
        "created_at": "Tue Sep 01 21:21:35 +0000 2009",
        "location": "Kingston New York",
        "follow_request_sent": null,
        "profile_link_color": "D02B55",
        "is_translator": false,
        "id_str": "70789458",
        "entities": {
          "url": {
            "urls": [
              {
                "expanded_url": null,
                "url": "",
                "indices": [
                  0,
                  0
                ]
              }
            ]
          },
          "description": {
            "urls": [
 
            ]
          }
        },
        "default_profile": false,
        "contributors_enabled": false,
        "favourites_count": 19,
        "url": null,
        "profile_image_url_https": "https://si0.twimg.com/profile_images/2219333930/Froggystyle_normal.png",
        "utc_offset": -18000,
        "id": 70789458,
        "profile_use_background_image": true,
        "listed_count": 1,
        "profile_text_color": "3E4415",
        "lang": "en",
        "followers_count": 63,
        "protected": false,
        "notifications": null,
        "profile_background_image_url_https": "https://si0.twimg.com/images/themes/theme5/bg.gif",
        "profile_background_color": "352726",
        "verified": false,
        "geo_enabled": false,
        "time_zone": "Eastern Time (US & Canada)",
        "description": "Science Fiction Writer, sort of. Likes Superheroes, Mole People, Alt. Timelines.",
        "default_profile_image": false,
        "profile_background_image_url": "https://a0.twimg.com/images/themes/theme5/bg.gif",
        "statuses_count": 1048,
        "friends_count": 63,
        "following": null,
        "show_all_inline_media": false,
        "screen_name": "MonkiesFist"
      },
      "in_reply_to_screen_name": null,
      "source": "web",
      "in_reply_to_status_id": null
    },
    {
      "coordinates": null,
      "favorited": false,
      "truncated": false,
      "created_at": "Fri Sep 21 22:51:18 +0000 2012",
      "id_str": "249279667666817024",
      "entities": {
        "urls": [
 
        ],
        "hashtags": [
          {
            "text": "freebandnames",
            "indices": [
              20,
              34
            ]
          }
        ],
        "user_mentions": [
 
        ]
      },
      "in_reply_to_user_id_str": null,
      "contributors": null,
      "text": "The Foolish Mortals #freebandnames",
      "metadata": {
        "iso_language_code": "en",
        "result_type": "recent"
      },
      "retweet_count": 0,
      "in_reply_to_status_id_str": null,
      "id": 249279667666817024,
      "geo": null,
      "retweeted": false,
      "in_reply_to_user_id": null,
      "place": null,
      "user": {
        "profile_sidebar_fill_color": "BFAC83",
        "profile_sidebar_border_color": "615A44",
        "profile_background_tile": true,
        "name": "Marty Elmer",
        "profile_image_url": "https://a0.twimg.com/profile_images/1629790393/shrinker_2000_trans_normal.png",
        "created_at": "Mon May 04 00:05:00 +0000 2009",
        "location": "Wisconsin, USA",
        "follow_request_sent": null,
        "profile_link_color": "3B2A26",
        "is_translator": false,
        "id_str": "37539828",
        "entities": {
          "url": {
            "urls": [
              {
                "expanded_url": null,
                "url": "https://www.omnitarian.me",
                "indices": [
                  0,
                  24
                ]
              }
            ]
          },
          "description": {
            "urls": [
 
            ]
          }
        },
        "default_profile": false,
        "contributors_enabled": false,
        "favourites_count": 647,
        "url": "https://www.omnitarian.me",
        "profile_image_url_https": "https://si0.twimg.com/profile_images/1629790393/shrinker_2000_trans_normal.png",
        "utc_offset": -21600,
        "id": 37539828,
        "profile_use_background_image": true,
        "listed_count": 52,
        "profile_text_color": "000000",
        "lang": "en",
        "followers_count": 608,
        "protected": false,
        "notifications": null,
        "profile_background_image_url_https": "https://si0.twimg.com/profile_background_images/106455659/rect6056-9.png",
        "profile_background_color": "EEE3C4",
        "verified": false,
        "geo_enabled": false,
        "time_zone": "Central Time (US & Canada)",
        "description": "Cartoonist, Illustrator, and T-Shirt connoisseur",
        "default_profile_image": false,
        "profile_background_image_url": "https://a0.twimg.com/profile_background_images/106455659/rect6056-9.png",
        "statuses_count": 3575,
        "friends_count": 249,
        "following": null,
        "show_all_inline_media": true,
        "screen_name": "Omnitarian"
      },
      "in_reply_to_screen_name": null,
      "source": "<a href=\"//twitter.com/download/iphone%5C%22\" rel=\"\\\"nofollow\\\"\">Twitter for iPhone</a>",
      "in_reply_to_status_id": null
    }
  ],
  "search_metadata": {
    "max_id": 250126199840518145,
    "since_id": 24012619984051000,
    "refresh_url": "?since_id=250126199840518145&q=%23freebandnames&result_type=mixed&include_entities=1",
    "next_results": "?max_id=249279667666817023&q=%23freebandnames&count=4&include_entities=1&result_type=mixed",
    "count": 4,
    "completed_in": 0.035,
    "since_id_str": "24012619984051000",
    "query": "%23freebandnames",
    "max_id_str": "250126199840518145"
  }
}`

func TestFastEncoding(t *testing.T) {

	var src *TwitterStruct
	{
		err := json.Unmarshal([]byte(TwitterJson), &src)
		if err != nil {
			panic(err)
		}
	}

	for j := 0; j < 5; j++ {
		N := 10000

		var (
			p1, p2 interface{}
		)

		start := time.Now()
		for i := 0; i < N; i++ {
			b, err := json.Marshal(src)
			if err != nil {
				panic(err)
			}
			var d1 TwitterStruct
			err = json.Unmarshal(b, &d1)
			if err != nil {
				panic(err)
			}
			if i == 0 {
				p1 = d1
			}
		}
		fmt.Println("std.Marshal", time.Since(start))

		start = time.Now()
		for i := 0; i < N; i++ {
			b, err := sonic.Marshal(src)
			if err != nil {
				panic(err)
			}
			var d1 TwitterStruct
			err = sonic.Unmarshal(b, &d1)
			if err != nil {
				panic(err)
			}
			if i == 0 {
				p1 = d1
			}
		}
		fmt.Println("sonic.Marshal", time.Since(start))

		start = time.Now()
		for i := 0; i < N; i++ {
			var d2 TwitterStruct
			err := cast.FAST.Convert(src, &d2)
			if err != nil {
				panic(err)
			}
			if i == 0 {
				p2 = d2
			}
		}
		fmt.Println("To", time.Since(start))

		{
			b1, _ := json.Marshal(p1)
			b2, _ := json.Marshal(p2)
			_ = json.Unmarshal(b1, &p1)
			_ = json.Unmarshal(b2, &p2)
		}

		if !reflect.DeepEqual(p1, p2) {
			t.Fatalf("\n%v\nnot equals\n%v\n", p1, p2)
		}
	}
}

func TestJsonEncoding(t *testing.T) {

	t.Run("", func(t *testing.T) {
		err := cast.JSON.Convert(make(chan int), nil)
		assert.Error(t, err, "json: unsupported type: chan int")
	})

	t.Run("", func(t *testing.T) {
		var a struct {
			Text string `json:"text"`
		}
		var b struct {
			Text *string `json:"text"`
		}
		a.Text = "hello, gopher!"
		err := cast.JSON.Convert(a, &b)
		assert.Nil(t, err)
		assert.Equal(t, *b.Text, a.Text)
	})
}
