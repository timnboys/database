package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MultiPanelTargets struct {
	*pgxpool.Pool
}

func newMultiPanelTargets(db *pgxpool.Pool) *MultiPanelTargets {
	return &MultiPanelTargets{
		db,
	}
}

func (p MultiPanelTargets) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS multi_panel_targets(
	"multi_panel_id" int4 NOT NULL,
	"panel_id" int8 NOT NULL,
	FOREIGN KEY("multi_panel_id") REFERENCES multi_panels("id") ON DELETE CASCADE,
	FOREIGN KEY ("panel_id") REFERENCES panels("message_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("multi_panel_id", "panel_id")
);
`
}

func (p *MultiPanelTargets) GetPanels(multiPanelId int) (panels []Panel, e error) {
	query := `
SELECT
	panels.message_id, panels.channel_id, panels.guild_id, panels.title, panels.content, panels.colour, panels.target_category, panels.reaction_emote, panels.welcome_message
FROM
	multi_panel_targets
INNER JOIN
	panels ON panels.message_id = multi_panel_targets.panel_id
WHERE
	"multi_panel_id" = $1
;`

	rows, err := p.Query(context.Background(), query, multiPanelId)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var panel Panel
		if err := rows.Scan(&panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.TargetCategory, &panel.ReactionEmote, &panel.WelcomeMessage); err != nil {
			e = err
			continue
		}

		panels = append(panels, panel)
	}

	return
}

func (p *MultiPanelTargets) GetMultiPanels(panelId uint64) (multiPanelIds []int, e error) {
	query := `
SELECT
	"multi_panel_id"
FROM
	multi_panel_targets
WHERE
	"panel_id" = $1
;`

	rows, err := p.Query(context.Background(), query, panelId)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			e = err
			continue
		}

		multiPanelIds = append(multiPanelIds, id)
	}

	return
}

func (p *MultiPanelTargets) Insert(multiPanelId int, panelId uint64) (err error) {
	query := `
INSERT INTO
	multi_panel_targets("multi_panel_id", "panel_id")
VALUES
	($1, $2) 
ON CONFLICT("multi_panel_id", "panel_id") DO
	NOTHING
;
`

	_, err = p.Exec(context.Background(), query, multiPanelId, panelId)
	return
}

func (p *MultiPanelTargets) DeleteAll(multiPanelId int) (err error) {
	query := `
DELETE FROM
	multi_panel_targets
WHERE
	"multi_panel_id"=$1
;`

	_, err = p.Exec(context.Background(), query, multiPanelId)
	return
}

func (p *MultiPanelTargets) Delete(multiPanelId int, panelId uint64) (err error) {
	query := `
DELETE FROM
	multi_panel_targets
WHERE
	"multi_panel_id"=$1
	AND
	"panel_id" = $2
;`

	_, err = p.Exec(context.Background(), query, multiPanelId, panelId)
	return
}
