$(document).ready(function() {
  var
    $headers     = $('body > h1'),
    $header      = $headers.first(),
    ignoreScroll = false,
    timer;

  $(window)
    .on('resize', function() {
      clearTimeout(timer);
      $headers.visibility('disable callbacks');

      $(document).scrollTop( $header.offset().top );

      timer = setTimeout(function() {
        $headers.visibility('enable callbacks');
      }, 500);
    });
  $headers
    .visibility({
      once: false,
      checkOnRefresh: true,
      onTopPassed: function() {
        $header = $(this);
      },
      onTopPassedReverse: function() {
        $header = $(this);
      }
    });
});

var SubmitForm = function() {
  $("#submit").addClass('disabled');
  var action  = $("#action").val();
  var message = $('#message').val();
  if (!message) {
    $("#submit").removeClass('disabled');
    $("#warning").text("Message is Empty").removeClass("hidden").addClass("visible");
    return false;
  }
  const data = {action, message};
  request(data, (res)=>{
    $("#result").text(res.message);
    $("#info").removeClass("hidden").addClass("visible");
  }, (e)=>{
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
    $("#submit").removeClass('disabled');
  });
};

var request = function(data, callback, onerror) {
  $.ajax({
    type:          'POST',
    dataType:      'json',
    contentType:   'application/json',
    scriptCharset: 'utf-8',
    data:          JSON.stringify(data),
    url:           {{ .Api }}
  })
  .done(function(res) {
    callback(res);
  })
  .fail(function(e) {
    onerror(e);
  });
};
