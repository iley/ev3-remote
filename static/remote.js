/*jslint browser: true */
/*jslint plusplus: true */
/*global WebSocket */

// TODO: Support multiple gamepads maybe?
(function () {
  'use strict';

  function round10(value, digits) {
      return Number(Math.round(value + 'e' + digits) + 'e-' + digits);
  }

  function TextConsole(container_id) {
    this.container = document.getElementById(container_id);
  }

  TextConsole.prototype.Write = function (text) {
    // TODO: escape HTML
    this.container.innerHTML += text + "<br>";
    this.container.scrollTop = this.container.scrollHeight;
  };

  function RemoteControl(websocket_url, text_console) {
    this.websocket_url = websocket_url;
    this.text_console = text_console;
    this.websocket = null;
    this.connected = false;
  }

  RemoteControl.prototype.Connect = function () {
    this.text_console.Write("Connecting to " + this.websocket_url);
    this.websocket = new WebSocket(this.websocket_url);
    this.websocket.onopen = function () {
      this.text_console.Write("Connected to websocket " + this.websocket_url);
      this.connected = true;
    }.bind(this);
    this.websocket.onclose = function () {
      this.text_console.Write("Websocket connection closed");
      this.connected = false;
    }.bind(this);
  };

  RemoteControl.prototype.Send = function (message) {
    if (!this.connected) {
      this.text_console.Write("Websocket not connected");
    } else {
      this.websocket.send(JSON.stringify(message));
    }
  };

  function RemoteApp(text_console, remote_control) {
    this.text_console = text_console;
    this.remote_control = remote_control;
    this.gamepad = null;
    this.ticking = false;
    this.axis_threshold = 0.05;
    this.axis_points = [0.25, 0.5, 0.75];
  }

  RemoteApp.prototype.Run = function () {
    this.remote_control.Connect();
    var gamepads = navigator.getGamepads();
    if (!gamepads) {
      this.text_console.Write("Gamepad support not avaliable");
      return;
    }
    this.WaitForGamepad();
  };

  RemoteApp.prototype.WaitForGamepad = function () {
    var gamepad = navigator.getGamepads()[0];
    if (!gamepad) {
      this.text_console.Write("Gamepad not found. Waiting...");
      window.setTimeout(this.WaitForGamepad.bind(this), 500);
    } else {
      this.text_console.Write("Gamepad found: " + gamepad.id);
      this.text_console.Write("Polling...");
      this.StartPolling();
    }
  };

  RemoteApp.prototype.StartPolling = function () {
    if (!this.ticking) {
      this.ticking = true;
      this.previous_timestamp = null;
      this.previous_axes = null;
      this.previous_button_values = null;
      this.Tick();
    }
  };

  RemoteApp.prototype.StopPolling = function () {
    this.ticking = false;
  };

  RemoteApp.prototype.Tick = function () {
    this.PollGamepadStatus();
    if (this.ticking) {
      window.requestAnimationFrame(this.Tick.bind(this));
    }
  };

  RemoteApp.prototype.PollGamepadStatus = function () {
    var i, axes, gamepad;
    gamepad = navigator.getGamepads() && navigator.getGamepads()[0];
    if (!gamepad) {
      this.text_console.Write("Gamepad disconnected");
      this.StopPolling();
      this.WaitForGamepad();
      return;
    }
    if (this.previous_timestamp === gamepad.timestamp) {
      return;  // Nothing changed.
    }
    this.UpdateButtons(gamepad.buttons);
    this.UpdateAxes(gamepad.axes);
    this.previous_timestamp = gamepad.timestamp;

  };

  // TODO: Support for analog buttons.
  RemoteApp.prototype.UpdateButtons = function (buttons) {
    var i, button_values = [];
    for (i = 0; i < buttons.length; ++i) {
      button_values[i] = buttons[i].pressed;
    }
    for (i = 0; i < button_values.length; ++i) {
      if (!this.previous_button_values ||
          this.previous_button_values[i] !== button_values[i]) {
        this.text_console.Write("Button #" + i + " " +
                                (button_values[i] ? "pressed" : "released"));
        this.remote_control.Send(
          {control: "button", id: i, value: button_values[i] ? 1 : 0}
        );
      }
    }
    this.previous_button_values = button_values;
  };

  RemoteApp.prototype.UpdateAxes = function (axes) {
    var i,
      normalized_axes = axes.slice();
    for (i = 0; i < axes.length; ++i) {
      normalized_axes[i] = round10(axes[i], 1);
      if (!this.previous_axes ||
          normalized_axes[i] !== this.previous_axes[i]) {
        this.text_console.Write("Axis #" + i + " = " + normalized_axes[i]);
        this.remote_control.Send(
          {control: "axis", id: i, value: normalized_axes[i]}
        );
      }
    }
    this.previous_axes = normalized_axes;
  };

  document.addEventListener("DOMContentLoaded", function () {
    var websocket_url = "ws://" + location.host + "/websocket",
      text_console = new TextConsole("mainConsole"),
      remote_control = new RemoteControl(websocket_url, text_console),
      app = new RemoteApp(text_console, remote_control);
    app.Run();
  });

}());
