(function() {
    'use strict';
    angular
        .module('DiggingApp')
        .directive('topicPassages', topicPassages);

    function topicPassages($location, $http, $log, webConfig, $routeParams) {
        var getWordsInTopic = function(scope) {
            scope.wordsInTopic = "No words yet...."
        }
        return {
            templateUrl: 'components/topicView/topicPassages.html',
            link: function(scope) {
                scope.debug = webConfig.debug;
                if (scope.debug) {
                    getWordsInTopic(scope);
                }
                scope.displayLimit = 50;
                var urlString = "/api/" + scope.main.dbActive + "/topic/" + $routeParams.topicID;
                var promise = $http.get(urlString);
                promise.then(function(response) {
                    scope.main.hideSearchForm = true;
                    scope.topicPassages = response.data;
                });
                scope.loadingData = false;
                scope.addMoreResults = function() {
                    scope.loadingData = true;
                    var lastWeight = scope.topicPassages[scope.topicPassages.length-1].topicWeight;
                    var update = $http.get(urlString + "?topicWeight=" + lastWeight);
                    update.then(function(response) {
                        Array.prototype.push.apply(scope.topicPassages, response.data);
                        scope.displayLimit += 100;
                        scope.loadingData = false;
                    });
                }
            }
        }
    }
})();
