var forsenE = (p) => {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get(p);
};
var loadVods = () => {
  let vods = [];
  let data = null;
  $.ajax({
    url: "/vods",
    type: "get",
    async: false,
    success: (response) => {
      data = response;
    },
  });
  data.vods.forEach((e) => {
    vods.push(e);
  });
  return vods;
};
var loadImageUrl = (data, emote) => {
  var emoteCompany = ["bttv", "7tv", "twitch", "ffz"];
  let emoteUrl = null;
  emoteCompany.forEach((element) => {
    data[element].forEach((e) => {
      if (e.code == emote) {
        emoteUrl = e.url;
      }
    });
  });
  if (emoteUrl == null) {
    return undefined;
  } else {
    return `<img src="${emoteUrl}"/>`;
  }
};
var loadEmotes = (logs) => {
  let data = null;
  $.ajax({
    url: "/emotes",
    type: "get",
    async: false,
    success: (response) => {
      data = response;
    },
  });
  for (const [key, value] of Object.entries(logs)) {
    value.forEach((e) => {
      if (e.emotes) {
        e.emotes.forEach(function (emote) {
          let verification = loadImageUrl(data, emote.code);
          if (verification == undefined) {
            data.twitch.push({ code: emote.code, url: emote.url });
          }
        });
      }
    });
  }
  return data;
};
var returnLinks = (link) => {
  return link.linkify({
    className: "externallink",
    rel: "nofollow noreferrer",
    target: "_blank",
  });
};
