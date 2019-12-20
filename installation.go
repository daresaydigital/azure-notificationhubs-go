package notificationhubs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
)

// Installation reads one specific installation
func (h *NotificationHub) Installation(ctx context.Context, installationID string) (raw []byte, installation *Installation, err error) {
	var (
		instURL = h.generateAPIURL(path.Join("installations", installationID))
	)

	raw, _, err = h.exec(ctx, getMethod, instURL, Headers{}, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, &installation)
	return
}

// Install sends a device installation to the Azure hub
func (h *NotificationHub) Install(ctx context.Context, installation Installation) (err error) {
	var (
		instURL = h.generateAPIURL(path.Join("installations", installation.InstallationID))
		headers = map[string]string{
			"Content-Type": "application/json",
		}
	)

	raw, err := json.Marshal(installation)
	if err != nil {
		return
	}

	_, _, err = h.exec(ctx, putMethod, instURL, headers, bytes.NewBuffer(raw))
	return
}

// Update sends a collection of installation changes to the Azure hub
func (h *NotificationHub) Update(ctx context.Context, installationID string, changes ...InstallationChange) (err error) {
	var (
		instURL = h.generateAPIURL(path.Join("installations", installationID))
		headers = map[string]string{
			"Content-Type": "application/json-patch+json",
		}
	)

	raw, err := json.Marshal(changes)
	if err != nil {
		return
	}

	_, _, err = h.exec(ctx, patchMethod, instURL, headers, bytes.NewBuffer(raw))
	return
}

// SetPushChannel sets the installation push channel
func SetPushChannel(pushChannel string) InstallationChange {
	return InstallationChange{Op: InstallationChangeReplace, Path: "/pushChannel", Value: pushChannel}
}

// SetTags sets the installation tags
func SetTags(tags ...string) InstallationChange {
	raw, _ := json.Marshal(tags)
	return InstallationChange{Op: InstallationChangeReplace, Path: "/tags", Value: string(raw)}
}

// AddTag adds a tag to the installation
func AddTag(tag string) InstallationChange {
	return InstallationChange{Op: InstallationChangeAdd, Path: "/tags", Value: tag}
}

// RemoveTag removes a tag from the installation
func RemoveTag(tag string) InstallationChange {
	return InstallationChange{Op: InstallationChangeRemove, Path: "/tags/" + tag}
}

// SetTemplates sets the installation templates
// Deprecated: doesn't appear to be supported
func SetTemplates(templates map[string]InstallationTemplate) InstallationChange {
	raw, _ := json.Marshal(templates)
	return InstallationChange{Op: InstallationChangeReplace, Path: "/templates", Value: string(raw)}
}

// AddTemplate adds a template to the installation
func AddTemplate(name string, template InstallationTemplate) InstallationChange {
	raw, _ := json.Marshal(template)
	return InstallationChange{Op: InstallationChangeAdd, Path: "/templates/" + name, Value: string(raw)}
}

// SetTemplateBody sets the body on a template in the installation
func SetTemplateBody(name, body string) InstallationChange {
	return InstallationChange{Op: InstallationChangeReplace, Path: fmt.Sprintf("/templates/%s/body", name), Value: body}
}

// SetTemplateHeaders sets the headers on a template in the installation
func SetTemplateHeaders(name string, headers map[string]string) InstallationChange {
	raw, _ := json.Marshal(headers)
	return InstallationChange{Op: InstallationChangeReplace, Path: fmt.Sprintf("/templates/%s/headers", name), Value: string(raw)}
}

// SetTemplateTags sets the tags on a template in the installation
func SetTemplateTags(name string, tags ...string) InstallationChange {
	raw, _ := json.Marshal(tags)
	return InstallationChange{Op: InstallationChangeReplace, Path: fmt.Sprintf("/templates/%s/tags", name), Value: string(raw)}
}

// AddTemplateTag adds a tag to a template in the installation
func AddTemplateTag(name, tag string) InstallationChange {
	return InstallationChange{Op: InstallationChangeAdd, Path: fmt.Sprintf("/templates/%s/tags", name), Value: tag}
}

// RemoveTemplateTag removes a tag from a template in the installation
func RemoveTemplateTag(name, tag string) InstallationChange {
	return InstallationChange{Op: InstallationChangeRemove, Path: fmt.Sprintf("/templates/%s/tags/%s", name, tag)}
}

// RemoveTemplate removes a template from the installation
func RemoveTemplate(name string) InstallationChange {
	return InstallationChange{Op: InstallationChangeRemove, Path: "/templates/" + name}
}

// SetSecondaryTiles sets the installation secondary tiles
// Deprecated: doesn't appear to be supported
func SetSecondaryTiles(secondaryTiles map[string]InstallationSecondaryTile) InstallationChange {
	raw, _ := json.Marshal(secondaryTiles)
	return InstallationChange{Op: InstallationChangeReplace, Path: "/secondaryTiles", Value: string(raw)}
}

// AddSecondaryTile adds a secondary tile to the installation
// Deprecated: doesn't appear to be supported
func AddSecondaryTile(name string, secondaryTile InstallationSecondaryTile) InstallationChange {
	raw, _ := json.Marshal(secondaryTile)
	return InstallationChange{Op: InstallationChangeAdd, Path: "/secondaryTiles/" + name, Value: string(raw)}
}

// SetSecondaryTilePushChannel sets the push channel on a secondary tile in the installation
func SetSecondaryTilePushChannel(name, pushChannel string) InstallationChange {
	return InstallationChange{Op: InstallationChangeReplace, Path: fmt.Sprintf("/secondaryTiles/%s/pushChannel", name), Value: pushChannel}
}

// SetSecondaryTileTags sets the tags on a secondary tile in the installation
func SetSecondaryTileTags(name string, tags ...string) InstallationChange {
	raw, _ := json.Marshal(tags)
	return InstallationChange{Op: InstallationChangeReplace, Path: fmt.Sprintf("/secondaryTiles/%s/tags", name), Value: string(raw)}
}

// AddSecondaryTileTag adds a tag to a secondary tile in the installation
func AddSecondaryTileTag(name, tag string) InstallationChange {
	return InstallationChange{Op: InstallationChangeAdd, Path: fmt.Sprintf("/secondaryTiles/%s/tags", name), Value: tag}
}

// RemoveSecondaryTileTag removes a tag from a secondary tile in the installation
func RemoveSecondaryTileTag(name, tag string) InstallationChange {
	return InstallationChange{Op: InstallationChangeRemove, Path: fmt.Sprintf("/secondaryTiles/%s/tags/%s", name, tag)}
}

// SetSecondaryTileTemplates sets the installation templates
func SetSecondaryTileTemplates(name string, templates map[string]InstallationTemplate) InstallationChange {
	raw, _ := json.Marshal(templates)
	return InstallationChange{Op: InstallationChangeReplace, Path: fmt.Sprintf("/secondaryTiles/%s/templates", name), Value: string(raw)}
}

// AddSecondaryTileTemplate adds a template to the installation
func AddSecondaryTileTemplate(name, templateName string, template InstallationTemplate) InstallationChange {
	raw, _ := json.Marshal(template)
	return InstallationChange{Op: InstallationChangeAdd, Path: fmt.Sprintf("/secondaryTiles/%s/templates/%s", name, templateName), Value: string(raw)}
}

// SetSecondaryTileTemplateBody sets the body on a template in the installation
func SetSecondaryTileTemplateBody(name, template, body string) InstallationChange {
	return InstallationChange{Op: InstallationChangeReplace, Path: fmt.Sprintf("/secondaryTiles/%s/templates/%s/body", name, template), Value: body}
}

// SetSecondaryTileTemplateHeaders sets the headers on a template in the installation
func SetSecondaryTileTemplateHeaders(name, template string, headers map[string]string) InstallationChange {
	raw, _ := json.Marshal(headers)
	return InstallationChange{Op: InstallationChangeReplace, Path: fmt.Sprintf("/secondaryTiles/%s/templates/%s/headers", name, template), Value: string(raw)}
}

// SetSecondaryTileTemplateTags sets the tags on a template in the installation
func SetSecondaryTileTemplateTags(name, template string, tags ...string) InstallationChange {
	raw, _ := json.Marshal(tags)
	return InstallationChange{Op: InstallationChangeReplace, Path: fmt.Sprintf("/secondaryTiles/%s/templates/%s/tags", name, template), Value: string(raw)}
}

// RemoveSecondaryTileTemplate removes a template from the installation
func RemoveSecondaryTileTemplate(name, template string) InstallationChange {
	return InstallationChange{Op: InstallationChangeRemove, Path: fmt.Sprintf("/secondaryTiles/%s/templates/%s", name, template)}
}

// RemoveSecondaryTile removes a secondary tile from the installation
// Deprecated: doesn't appear to be supported
func RemoveSecondaryTile(name string) InstallationChange {
	return InstallationChange{Op: InstallationChangeRemove, Path: "/secondaryTiles/" + name}
}

// Uninstall sends a device installation delete to the Azure hub
func (h *NotificationHub) Uninstall(ctx context.Context, installationID string) (err error) {
	var (
		instURL = h.generateAPIURL(path.Join("installations", installationID))
		headers = map[string]string{
			"Content-Type": "application/json",
		}
	)

	_, _, err = h.exec(ctx, deleteMethod, instURL, headers, nil)
	return
}
