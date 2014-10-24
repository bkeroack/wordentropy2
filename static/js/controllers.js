var wordentropy2App = angular.module('wordentropy2', []);

wordentropy2App.controller('PassphrasesController', function ($scope, $http, $location) {
	var home = $location.protocol() + "://" + $location.host() + ":" + $location.port();

	$scope.length = 5;
	$scope.count = 5;

	$scope.passphrases = [];

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
						"strength_label": results.score > 3 ? "strong" : "weak"
					}
					$scope.passphrases.push(phrase_obj);
				}
			}
		}).error(function(data, status, headers, config) {
			console.log("Error getting passphrases: " + data);
		});
	};

});