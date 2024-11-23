package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
			{
				"id": "_pb_users_auth_",
				"created": "2024-11-20 18:00:17.871Z",
				"updated": "2024-11-20 18:16:25.781Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null,
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					},
					{
						"system": false,
						"id": "xolbud0n",
						"name": "bookmarks",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "r46wuuwnllboa9m",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "cwbs2kky",
						"name": "tags",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "c0x0ngay68c3wjc",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": null
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"onlyVerified": false,
					"requireEmail": false
				}
			},
			{
				"id": "r46wuuwnllboa9m",
				"created": "2024-11-20 18:08:07.175Z",
				"updated": "2024-11-20 18:41:59.960Z",
				"name": "bookmarks",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "wieyqovr",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "l7hw26k7",
						"name": "url",
						"type": "url",
						"required": true,
						"presentable": true,
						"unique": false,
						"options": {
							"exceptDomains": [],
							"onlyDomains": []
						}
					},
					{
						"system": false,
						"id": "g6ijsrvu",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": true,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "gfmqxf16",
						"name": "tags",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "c0x0ngay68c3wjc",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": null
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_9xlgK9L` + "`" + ` ON ` + "`" + `bookmarks` + "`" + ` (` + "`" + `url` + "`" + `)"
				],
				"listRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"viewRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"createRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"updateRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"deleteRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"options": {}
			},
			{
				"id": "c0x0ngay68c3wjc",
				"created": "2024-11-20 18:10:31.210Z",
				"updated": "2024-11-20 18:37:43.175Z",
				"name": "tags",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "agecscb8",
						"name": "name",
						"type": "text",
						"required": true,
						"presentable": true,
						"unique": false,
						"options": {
							"min": 1,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "umwx41ke",
						"name": "color",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 4,
							"max": null,
							"pattern": "^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
						}
					},
					{
						"system": false,
						"id": "gseqdwlk",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": true,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "lcr1inpk",
						"name": "bookmarks",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "r46wuuwnllboa9m",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": null
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_r7giYgS` + "`" + ` ON ` + "`" + `tags` + "`" + ` (` + "`" + `name` + "`" + `)"
				],
				"listRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"viewRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"createRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"updateRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"deleteRule": "@request.auth.id != \"\" && user = @request.auth.id ",
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
