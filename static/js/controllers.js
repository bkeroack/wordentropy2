var wordentropy2App = angular.module('wordentropy2', []);

wordentropy2App.controller('PassphrasesController', function ($scope, $http, $location) {
	var home = $location.protocol() + "://" + $location.host() + ":" + $location.port();

	$scope.length = 4;
	$scope.count = 5;

	$scope.passphrases = [];

	$scope.resetError = function() {
		$scope.error_alert = false;
		$scope.error_alert_msg = "";
		$scope.error_alert_class = "";
	};

	$scope.getPassphrases = function() {
		
		var url = home + "/passphrases?length=" + encodeURIComponent($scope.length) + "&count="
		+ encodeURIComponent($scope.count);
		
		$http.get(url).success(function(data, status, headers, config) {
			if (data.hasOwnProperty("Passphrases") && typeof(zxcvbn) == 'function') {
				$scope.passphrases = [];
				for (var i in data.Passphrases) {
					var phrase = data.Passphrases[i];
					var results = zxcvbn(phrase);
					var phrase_obj = {
						"phrase": phrase,
						"bits": results.entropy,
						"crack_time": results.crack_time_display,
						"strength": results.score,
						"strength_label": results.score > 3 ? "strong" : (results.score > 2 ? "decent" : "weak"),
						"strength_class": results.score > 3 ? "alert-success" : (results.score > 2 ? "alert-warning" : "alert-danger")
					}
					$scope.passphrases.push(phrase_obj);
				}
				$scope.resetError();
			}
		}).error(function(data, status, headers, config) {
			var error_msg = "Error getting passphrases: " + JSON.stringify(data); 
			console.log(error_msg);
			$scope.error_alert_msg = error_msg;
			$scope.error_alert_class = "alert-danger";
			$scope.error_alert = true;
		});
	};

	$scope.resetError();

});