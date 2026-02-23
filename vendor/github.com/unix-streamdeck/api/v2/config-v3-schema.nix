{ lib, ... }:

with lib;

let
 obsConnectionInfoV2Type = types.submodule {
    options = {
      host = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "OBS WebSocket host";
      };
      port = mkOption {
        type = types.nullOr types.int;
        default = null;
        description = "OBS WebSocket port";
      };
    };
  };

  knobActionV3Type = types.submodule {
    options = {
      switch_page = mkOption {
        type = types.nullOr types.int;
        default = null;
        description = "Page number to switch to";
      };
      keybind = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Keybind to trigger";
      };
      command = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Command to execute";
      };
      brightness = mkOption {
        type = types.nullOr types.int;
        default = null;
        description = "Brightness level";
      };
      url = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "URL to open";
      };
      obs_command = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "OBS command to execute";
      };
      obs_command_params = mkOption {
        type = types.nullOr (types.attrsOf types.str);
        default = null;
        description = "OBS command parameters";
      };
    };
  };

  keyConfigV3Type = types.submodule {
    options = {
      icon = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Path to icon file";
      };
      switch_page = mkOption {
        type = types.nullOr types.int;
        default = null;
        description = "Page number to switch to";
      };
      text = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Text to display";
      };
      text_size = mkOption {
        type = types.nullOr types.int;
        default = null;
        description = "Text size";
      };
      text_alignment = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Text alignment (left, center, right)";
      };
      keybind = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Keybind to trigger";
      };
      command = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Command to execute";
      };
      brightness = mkOption {
        type = types.nullOr types.int;
        default = null;
        description = "Brightness level";
      };
      url = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "URL to open";
      };
      key_hold = mkOption {
        type = types.nullOr types.int;
        default = null;
        description = "Key hold duration in milliseconds";
      };
      obs_command = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "OBS command to execute";
      };
      obs_command_params = mkOption {
        type = types.nullOr (types.attrsOf types.str);
        default = null;
        description = "OBS command parameters";
      };
      icon_handler = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Icon handler name";
      };
      key_handler = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Key handler name";
      };
      icon_handler_fields = mkOption {
        type = types.nullOr types.attrs;
        default = null;
        description = "Icon handler configuration fields";
      };
      key_handler_fields = mkOption {
        type = types.nullOr types.attrs;
        default = null;
        description = "Key handler configuration fields";
      };
      shared_handler_fields = mkOption {
        type = types.nullOr types.attrs;
        default = null;
        description = "Shared handler configuration fields";
      };
    };
  };

  knobConfigV3Type = types.submodule {
    options = {
      icon = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Path to icon file";
      };
      text = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Text to display";
      };
      text_size = mkOption {
        type = types.nullOr types.int;
        default = null;
        description = "Text size";
      };
      text_alignment = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Text alignment (left, center, right)";
      };
      lcd_handler = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "LCD handler name";
      };
      knob_or_touch_handler = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = "Knob or touch handler name";
      };
      lcd_handler_fields = mkOption {
        type = types.nullOr types.attrs;
        default = null;
        description = "LCD handler configuration fields";
      };
      knob_or_touch_handler_fields = mkOption {
        type = types.nullOr types.attrs;
        default = null;
        description = "Knob or touch handler configuration fields";
      };
      shared_handler_fields = mkOption {
        type = types.nullOr types.attrs;
        default = null;
        description = "Shared handler configuration fields";
      };
      knob_press_action = mkOption {
        type = types.nullOr knobActionV3Type;
        default = null;
        description = "Action on knob press";
      };
      knob_turn_up_action = mkOption {
        type = types.nullOr knobActionV3Type;
        default = null;
        description = "Action on knob turn up";
      };
      knob_turn_down_action = mkOption {
        type = types.nullOr knobActionV3Type;
        default = null;
        description = "Action on knob turn down";
      };
    };
  };

  keyV3Type = types.submodule {
    options = {
      application = mkOption {
        type = types.nullOr (types.attrsOf keyConfigV3Type);
        default = null;
        description = "Application-specific key configurations";
      };
    };
  };

  knobV3Type = types.submodule {
    options = {
      application = mkOption {
        type = types.nullOr (types.attrsOf knobConfigV3Type);
        default = null;
        description = "Application-specific knob configurations";
      };
    };
  };

  pageV3Type = types.submodule {
    options = {
      keys = mkOption {
        type = types.listOf keyV3Type;
        default = [];
        description = "List of keys on this page";
      };
      knobs = mkOption {
        type = types.listOf knobV3Type;
        default = [];
        description = "List of knobs on this page";
      };
    };
  };

  deckV3Type = types.submodule {
    options = {
      serial = mkOption {
        type = types.str;
        description = "Serial number of the Stream Deck";
      };
      pages = mkOption {
        type = types.listOf pageV3Type;
        default = [];
        description = "List of pages for this deck";
      };
    };
  };

  configV3Type = types.submodule {
    options = {
      modules = mkOption {
        type = types.nullOr (types.listOf types.str);
        default = null;
        description = "List of modules to load";
      };
      decks = mkOption {
        type = types.listOf deckV3Type;
        default = [];
        description = "List of Stream Deck configurations";
      };
      obs_connection_info = mkOption {
        type = types.nullOr obsConnectionInfoV2Type;
        default = null
        description = "OBS WebSocket connection information";
      };
    };
  };

in {
  type = configV3Type;
}
