$(document).ready(function () {
  var time = 0;
  var begin = forsenE("from");
  var end = forsenE("to");
  var vod_id = forsenE("vod");
  var allVODs = [];
  var page = 1;
  $("#player").css("display", "flex");
  if (end && begin && vod_id) {
    if (end != "" && begin != "" && vod_id != "") {
      $("#player").show();
      $("#browse").hide();
      var replacedVideo = document.createElement("video");
      replacedVideo.controls = true;
      replacedVideo.autoplay = true;
      replacedVideo.muted = true;
      replacedVideo.style.width = "100%";
      replacedVideo.style.objectFit = "contain";
      replacedVideo.style.height = "100%";
      document.querySelector("#video-player").appendChild(replacedVideo);
      var videoSrc = `https://surskity.ams1.vultrobjects.com/vods/${vod_id}.mp4`;
      replacedVideo.src = videoSrc;
      replacedVideo.currentTime = time;
      replacedVideo.pause();
      var forsen = new forsenPossessed(begin, end, replacedVideo);
      replacedVideo.addEventListener("play", function () {
        forsen.startPepeLaughing();
      });
      replacedVideo.addEventListener("pause", function () {
        forsen.stopPepeLaughing();
      });
      $("#dec-delay-button").click(function () {
        delay = Number($("#delay").val()) - 1;
        $("#delay").val(delay);
      });

      $("#inc-delay-button").click(function () {
        delay = Number($("#delay").val()) + 1;
        $("#delay").val(delay);
      });
      $("body").css("overflow", "hidden");
    }
  } else {
    $("#player").hide();
    $("#browse").show();
    $("#page-number").text(page);
    $("#vod-list").empty();
    let r = loadVods();
    if (r.length > 12) {
      $("#next-page-button").removeClass("disabled");
    }
    allVODs = r;
    r.slice(0, 12).forEach((vod) => {
      start = moment(vod.from.replace("+00:00", "Z")).utc();
      finish = moment(vod.to.replace("+00:00", "Z")).utc();
      let difference = moment.duration(finish.diff(start));
      if (difference.hours() == 0 && difference.days() == 0) {
        var formatted = `${String(difference.minutes()).padStart(
          2,
          "0"
        )}:${String(difference.seconds()).padStart(2, "0")}`;
      } else {
        var formatted = `${
          difference.hours() + difference.days() * 24
        }:${String(difference.minutes()).padStart(2, "0")}:${String(
          difference.seconds()
        ).padStart(2, "0")}`;
      }
      $("#vod-list").append(
        '<div id="vod-item"><a href="?vod=' +
          vod.vod_id +
          "&from=" +
          vod.from +
          "&to=" +
          vod.to +
          '"><div class="image-container"><img src="img/kys.png"></div></a><div class="info-container" style="padding-left: 10px;"><div class="info">' +
          start.format("MMMM") +
          " " +
          vod.from.split("T")[0].split("-")[2] +
          ", " +
          vod.from.split("T")[0].split("-")[0] +
          '</div><div class="info">' +
          formatted +
          "</div></div></div>"
      );
    });
  }
  if (allVODs) {
    if (allVODs.length > 0) {
      $("#next-page-button").click(function () {
        if (page != Math.ceil(allVODs.length / 12)) {
          page += 1;
          $("#page-number").text(page);
          $("#vod-list").empty();
          allVODs.slice((page - 1) * 12, page * 12).forEach((fors) => {
            start = moment(fors.from.replace("+00:00", "Z")).utc();
            finish = moment(fors.to.replace("+00:00", "Z")).utc();
            let difference = moment.duration(finish.diff(start));
            if (difference.hours() == 0 && difference.days() == 0) {
              var formatted = `${String(difference.minutes()).padStart(
                2,
                "0"
              )}:${String(difference.seconds()).padStart(2, "0")}`;
            } else {
              var formatted = `${
                difference.hours() + difference.days() * 24
              }:${String(difference.minutes()).padStart(2, "0")}:${String(
                difference.seconds()
              ).padStart(2, "0")}`;
            }
            $("#vod-list").append(
              '<div id="vod-item"><a href="?vod=' +
                fors.vod_id +
                "&from=" +
                fors.from +
                "&to=" +
                fors.to +
                '"><div class="image-container"><img src="img/kys.png"></div></a><div class="info-container" style="padding-left: 10px;"><div class="info">' +
                start.format("MMMM") +
                " " +
                fors.from.split("T")[0].split("-")[2] +
                ", " +
                fors.from.split("T")[0].split("-")[0] +
                '</div><div class="info">' +
                formatted +
                "</div></div></div>"
            );
          });
        }

        if (page === Math.ceil(allVODs.length / 12)) {
          $("#next-page-button").addClass("disabled");
        } else {
          $("#next-page-button").removeClass("disabled");
        }

        if (page === 1) {
          $("#previous-page-button").addClass("disabled");
        } else {
          $("#previous-page-button").removeClass("disabled");
        }
      });
      $("#previous-page-button").click(function () {
        if (page > 1) {
          page -= 1;
          $("#page-number").text(page);
          $("#vod-list").empty();
          allVODs.slice((page - 1) * 12, page * 12).forEach((fors) => {
            start = moment(fors.from.replace("+00:00", "Z")).utc();
            finish = moment(fors.to.replace("+00:00", "Z")).utc();
            let difference = moment.duration(finish.diff(start));
            if (difference.hours() == 0 && difference.days() == 0) {
              var formatted = `${String(difference.minutes()).padStart(
                2,
                "0"
              )}:${String(difference.seconds()).padStart(2, "0")}`;
            } else {
              var formatted = `${
                difference.hours() + difference.days() * 24
              }:${String(difference.minutes()).padStart(2, "0")}:${String(
                difference.seconds()
              ).padStart(2, "0")}`;
            }
            $("#vod-list").append(
              '<div id="vod-item"><a href="?vod=' +
                fors.vod_id +
                "&from=" +
                fors.from +
                "&to=" +
                fors.to +
                '"><div class="image-container"><img src="img/kys.png"></div></a><div class="info-container" style="padding-left: 10px;"><div class="info">' +
                start.format("MMMM") +
                " " +
                fors.from.split("T")[0].split("-")[2] +
                ", " +
                fors.from.split("T")[0].split("-")[0] +
                '</div><div class="info">' +
                formatted +
                "</div></div></div>"
            );
          });
          if (page === Math.ceil(allVODs.length / 12)) {
            $("#next-page-button").addClass("disabled");
          } else {
            $("#next-page-button").removeClass("disabled");
          }

          if (page === 1) {
            $("#previous-page-button").addClass("disabled");
          } else {
            $("#previous-page-button").removeClass("disabled");
          }
        }
      });
    }
  }
});
