var wordentropy2App = angular.module('wordentropy2', []);

wordentropy2App.controller('PassphrasesController', function ($scope, $http, $location) {
	var home = $location.protocol() + "://" + $location.host() + ":" + $location.port();

	$scope.length = 4;
	$scope.count = 5;

	$scope.passphrases = [];
	$scope.examples = [];
	$scope.prudish = false;
	$scope.no_spaces = false;
	$scope.add_digit = false;
	$scope.add_symbol = false;


	$scope.resetError = function() {
		$scope.error_alert = false;
		$scope.error_alert_msg = "";
		$scope.error_alert_class = "";
	};

	$scope.getRandomString = function(length) {
		var random_string = "";
	    var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

	    for (var i=0; i < length; i++) {
	        random_string += chars.charAt(Math.floor(Math.random() * chars.length));
	    }
	    return random_string;
	};

	// Get a word from the generated passphrases of length 4-8
	$scope.getWord = function() {
		for (var i in $scope.passphrases) {
			var words = $scope.passphrases[i].phrase.split(" ");
			for (var j in words) {
				if (words[j].length >= 4 && words[j].length <= 8) {
					return words[j]
				}
			}
		}
		// none found
		var fallback = ["hacker", "password", "foobar", "secure"];
		return fallback[Math.floor(Math.random()*fallback.length)];
	};

	$scope.substituteWord = function(word) {
		var substitutions = {
			'a': '4',
			'e': '3',
			'o': '0',
			'l': '1',
			's': '5',
			't': '7'
		};
		var special_chars = ['!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '+', '='];
		var new_word = "";
		var r = 0;

		for (var i in word) {
			var letter = word[i]
			if (letter in substitutions) {
				if (Math.random() > 0.5) {
					new_word += substitutions[letter];
					r += 1;
				} else {
					new_word += letter;
				}
			} else {
				new_word += letter;
			}
		}

		// require at least one subsitution
		if (r == 0) {
			new_word = "";
			for (i in word) {
				letter = word[i];
				if (letter in substitutions && r == 0) {
					new_word += substitutions[letter];
					r += 1;
				} else {
					new_word += letter;
				}
			}
		}

		// append a special character
		new_word += special_chars[Math.floor(Math.random()*special_chars.length)];
		return new_word;
	};

	$scope.getAvgEntropy = function() {
		var sum = 0;
		for (var i in $scope.passphrases) {
			sum += $scope.passphrases[i].bits;
		}
		return sum / $scope.passphrases.length;
	};

	$scope.generateExamples = function() {
		// 3 examples: word with substitutions, word with number, random string
		$scope.examples = [];

		var avg_ent = $scope.getAvgEntropy();

		var word = $scope.getWord();
		var ex1 = word;
		ex1 += Math.floor(Math.random()*10);
		ex1 += Math.floor(Math.random()*10);
		var results = zxcvbn(ex1);
		$scope.examples.push({
			"password": ex1,
			"description": "word with numbers",
			"mem_comparison": "",
			"sec_comparison": results.entropy < avg_ent ? "much less" : "more",
			"bits": results.entropy,
			"crack_time": results.crack_time_display,
			"strength": results.score,
			"strength_label": results.score > 3 ? "strong" : (results.score > 2 ? "decent" : "weak"),
			"strength_class": results.score > 3 ? "alert-success" : (results.score > 2 ? "alert-warning" : "alert-danger")
		});

		var ex2 = $scope.substituteWord(word);
		results = zxcvbn(ex2);
		$scope.examples.push({
			"password": ex2,
			"description": "substituted word with symbol",
			"mem_comparison": "harder to remember and",
			"sec_comparison": results.entropy < avg_ent ? "less" : "more",
			"bits": results.entropy,
			"crack_time": results.crack_time_display,
			"strength": results.score,
			"strength_label": results.score > 3 ? "strong" : (results.score > 2 ? "decent" : "weak"),
			"strength_class": results.score > 3 ? "alert-success" : (results.score > 2 ? "alert-warning" : "alert-danger")
		});

		var ex3 = $scope.getRandomString(8);
		results = zxcvbn(ex3);
		$scope.examples.push({
			"password": ex3,
			"description": "random string",
			"mem_comparison": "harder to remember and",
			"sec_comparison": results.entropy < avg_ent ? "less" : "more",
			"bits": results.entropy,
			"crack_time": results.crack_time_display,
			"strength": results.score,
			"strength_label": results.score > 3 ? "strong" : (results.score > 2 ? "decent" : "weak"),
			"strength_class": results.score > 3 ? "alert-success" : (results.score > 2 ? "alert-warning" : "alert-danger")
		});
	};

	$scope.getPassphrases = function() {
		
		var url = home + "/passphrases?length=" + encodeURIComponent($scope.length) + "&count="
		+ encodeURIComponent($scope.count);

		var options = ["prudish", "no_spaces", "add_digit", "add_symbol"];
		for (var i in options) {
			if ($scope[options[i]]) {
				url += "&" + options[i] + "=" + "true";
			}
		}
		
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
				$scope.generateExamples();
				$scope.resetError();
			}
		}).error(function(data, status, headers, config) {
			var error_msg = "Error getting passphrases: " + JSON.stringify(data); 
			console.log(error_msg);
			$scope.error_alert_msg = error_msg;
			$scope.error_alert_class = "alert-danger";
			$scope.error_alert = true;
			$scope.passphrases = [];
			$scope.examples = [];
		});
	};

	$scope.resetError();

});