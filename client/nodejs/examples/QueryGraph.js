/**
 * Query for members in the House of Bonaparte graph 
 * born between the start and end years
 * and display the member attributes.
 * 
 * Usage : node QueryGraph.js -startyear 1900 -endyear 2000
 * 
 */
var TGException  = require('../lib/exception/TGException').TGException,
    conFactory   = require('../lib/connection/TGConnectionFactory'),
	PrintUtility = require('../lib/utils/PrintUtility'),
	TGLogManager = require('../lib/log/TGLogManager'),
	TGLogLevel   = require('../lib/log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function main() {
	
	var startYear = 1700;
	var endYear = 1800;

	var startYearIndex = -1;
	var endYearIndex = -1;
	process.argv.forEach(function (val, index, array) {
		console.log(index + ': ' + val); 
		if (val==='-startyear') {
			startYearIndex = index+1;
		}
		
		if (val==='-endyear') {
			endYearIndex = index+1;
		}

		if(index===startYearIndex) {
			startYear = val;
		}
		
		if(index===endYearIndex) {
			endYear = val;
		}
	});

	
	var logger = TGLogManager.getLogger();
	logger.setLevel(TGLogLevel.Info);
			
	var url = "tcp://192.168.1.6:8222";
	var user = "napoleon";
	var pwd = "bonaparte";
	var conn = conFactory.getFactory().createConnection(url, user, pwd, null);
	conn.connect(function(){
		var gof = conn.getGraphObjectFactory();
		if (!gof) {
			throw new TGException("Graph object not found");
		}

		conn.getGraphMetadata(true, function(gmd){
			logger.logInfo("Querying for member born between %d and %d", startYear, endYear);
	        var queryString = "@nodetype = 'houseMemberType' and yearBorn > " + startYear + " and yearBorn < " + endYear + ";";
	        conn.executeQuery(queryString, null, function(resultSet){
	          	if (resultSet !== null) {
	                while (resultSet.hasNext()) {
	                    var houseMember = resultSet.next();
	                    
	                    logger.logInfo("House member '%s' found",houseMember.getAttribute("memberName").getValue());
	                    houseMember.getAttributes().forEach(function(attr) {
	                        if (attr.getValue() === null) {
	                         	logger.logInfo("\t%s: %s", attr.getName(), "");
	                        }
	                        else {
	                        	logger.logInfo("\t%s: %s", attr.getName(), attr.getValue());
	                        }
	                    });
	                }
	            } else {
	          		logger.logInfo("Querying for member born between %d and %d not found", startYear, endYear);
	          	}
	          		
	    		if (conn !== null) {
	    			conn.disconnect();
	    		}
	        });
		});
	});
}

main();