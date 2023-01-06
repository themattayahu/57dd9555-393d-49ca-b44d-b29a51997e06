var forsenPossessed = function (begin, end, player) {
  this.previousTimeOffset = -1;
  this.start = begin.replace(" ", "+").replace("+00:00", "Z");
  this.finish = end.replace(" ", "+").replace("+00:00", "Z");
  this.begin = moment(begin.replace(" ", "+").replace("+00:00", "Z")).utc();
  this.end = moment(end.replace(" ", "+").replace("+00:00", "Z")).utc();
  this.chat = null;
  this.emotes = null;
  this.videoPlayer = player;
  this.delay = $("#delay");
  this.status = "loading";
  this.chatStream = $("#chat-stream");
  this.actualPreviousTimeOffset = 0;
  this.bottomDetector = true;
  var self = this;
  $.ajax({
    url: "/logs",
    type: "get",
    async: false,
    data: {
      from: self.start,
      to: self.finish,
    },
    success: (data) => {
      self.chat = data;
    },
  });
  self.emoteCompany = ["bttv", "7tv", "twitch", "ffz"];
  self.emotes = loadEmotes(self.chat);
  self.emoteList = [];
  self.emoteCompany.forEach((element) => {
    self.emotes[element].forEach((e) => {
      self.emoteList.push(e.code);
    });
  });
  self.startPepeLaughing = function () {
    self.status = "running";
  };
  self.stopPepeLaughing = function () {
    self.status = "paused";
  };
  self.slaveFunction = () => {
    if (self.status == "running" && self.chat) {
      var currentTimeOffset = Math.floor(self.videoPlayer.currentTime);
      var utcFormat = [];
      var timestamps = [];
      if (currentTimeOffset != self.previousTimeOffset) {
        timeDifference = currentTimeOffset - self.actualPreviousTimeOffset;

        timestamps.push(
          self.begin.clone().unix() -
            Number(self.delay.val()) +
            currentTimeOffset
        );

        if (timeDifference > 1 && timeDifference < 30) {
          for (let i = 1; i < timeDifference; i++) {
            timestamps.push(
              self.begin.clone().unix() -
                Number(self.delay.val()) +
                currentTimeOffset -
                i
            );
          }
        }
        timestamps.forEach((element) => {
          if (element in self.chat) {
            utcFormat.unshift(element);
          }
        });
        let msgAmount = 0;
        for (let i = 0; i < utcFormat.length; i++) {
          msgAmount += self.chat[utcFormat[i]].length;
        }

        var randomTimeouts = Array.from({ length: msgAmount }, () =>
          Math.random()
        );
        randomTimeouts.sort((a, b) => a - b);

        i = 0;

        utcFormat.forEach((element) => {
          self.chat[element].forEach(function (chatLine) {
            setTimeout(function () {
              let formatedMessage = chatLine.message;
	      formatedMessage = returnLinks(formatedMessage);
              let messageArray = formatedMessage.split(" ");
              for (let messageItem of messageArray) {
                for (let emote of self.emoteList) {
                  if (messageItem == emote) {
                    let emoteHTML = loadImageUrl(self.emotes, emote);
                    formatedMessage = formatedMessage.replace(emote, emoteHTML);
                  }
                }
              }
	      let usernameFormatted =
                '<span class="user">' + chatLine.username + "</span>";
              if (chatLine.username == "bot") {
                usernameFormatted =
                  '<span class="user-bot">' + chatLine.username + "</span>";
              }
              self.chatStream.append(
                '<div class="chat-msg">' +
                  '<span class="time">' +
                  moment.unix(element).utc().format("HH:mm:ss") +
                  "</span>" + usernameFormatted + "<span>:</span> " +
                  '<span class="message">' +
                  formatedMessage +
                  "</span>" +
                  "</div>"
              );
              if (self.bottomDetector) {
                self.chatStream.scrollTop(self.chatStream[0].scrollHeight);
              }
            }, randomTimeouts[i] * (1000 / 1 - 25));
            i++;
          });
        });

        self.actualPreviousTimeOffset = currentTimeOffset;
        if (300 !== 0 && 300 !== "0" && 300 !== "") {
          if (self.chatStream.children().length > 300) {
            removeLine =
              "#chat-stream div:lt(" +
              (self.chatStream.children().length - 300) +
              ")";
            $(removeLine).remove();
          }
        }
      }

      self.previousTimeOffset = currentTimeOffset;
    }
  };
  self.videoPlayer.addEventListener("progress", function () {
    self.actualPreviousTimeOffset = Math.floor(self.videoPlayer.currentTime);
  });
  $("#chat-stream").on("scroll", function () {
    if (
      self.chatStream.scrollTop() + self.chatStream.innerHeight() >=
      self.chatStream[0].scrollHeight - 80
    ) {
      self.bottomDetector = true;
    } else {
      self.bottomDetector = false;
    }
  });

  self.slaveStarts = () => {
    self.slaveFunction();
    self.slaveInterval = window.setTimeout(self.slaveStarts, 1000 / 1);
  };
  self.slaveStarts();
};
