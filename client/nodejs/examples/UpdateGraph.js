/**
 * For a given member of the House, update the attributes
 * 
 * Usage : java UpdateGraph [options]
 * 
 *  where options are:
 *   -memberName <memberName> Required. Member name - "Napoleon Bonaparte"
 *   -crownName <crownName>   Optional. Name while reigning - "Napoleon XVIII"
 *   -crownTitle <crownTitle> Optional. Title while reigning - "King of USA"    
 *   -houseHead <houseHead>   Optional. Head of the house - true or false
 *   -yearBorn <yearBorn>     Optional. Year of birth - 2004
 *   -yearDied <yearDied>     Optional. Year of death - 2016 or null if still alive
 *   -reignStart <reignStart> Optional. Date reign starts (format dd MMM yyyy) - 20 Jan 2008 or null if never reigned
 *   -reignEnd <reignEnd>     Optional. Date reign ends (format dd MMM yyyy) - 08 Nov 2016 or null if never reigned or still reigning
 *   
 *  For instance to update the house member named "Napoleon Bonaparte" :
 *  node UpdateGraph.js -memberName "Napoleon Bonaparte" -crownName "Napoleon XVIII" -crownTitle "King of USA" -yearDied null -reignEnd "31 Jan 2016"
 *
 */
var TGException  = require('../lib/exception/TGException').TGException,
    conFactory   = require('../lib/connection/TGConnectionFactory'),
	PrintUtility = require('../lib/utils/PrintUtility'),
	TGLogManager = require('../lib/log/TGLogManager'),
	TGLogLevel   = require('../lib/log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();
logger.setLevel(TGLogLevel.Info);
	
function main() {
	var url = "tcp://192.168.1.6:8222";
	var user = "napoleon";
	var pwd = "bonaparte";

	var memberName = "Napoleon Bonaparte";
	var crownName = "Grand Napoleon";
	var crownTitle = "King of the world";
	var yearBorn = null;
	var yearDied = "2016";
	var reignStart = "8 Nov 2001";
	var reignEnd= null;
	var houseHead = null;

	var memberNameIndex = -1;
	var crownNameIndex = -1;
	var crownTitleIndex = -1;
	var yearBornIndex = -1;
	var yearDiedIndex = -1;
	var reignStartIndex = -1;
	var reignEndIndex= -1;
	var houseHeadIndex = -1;
	process.argv.forEach(function (val, index, array) {
		if (val==="-memberName") {
			memberNameIndex = index+1;
		} else if (val==="-crownName") {
			crownNameIndex = index+1;
		} else if (val==="-crownTitle") {
			crownTitleIndex = index+1;
		} else if (val==="-yearBorn") {
			yearBornIndex = index+1;
		} else if (val==="-yearDied") {
			yearDiedIndex = index+1;
		} else if (val==="-houseHead") {
			houseHeadIndex = index+1;
		} else if (val==="-reignStart") {
			reignStartIndex = index+1;
		} else if (val==="-reignEnd") {
			reignEndIndex = index+1;
		}
		
		if (index===memberNameIndex) {
			memberName = val;
		} else if (index===crownNameIndex) {
			crownName = val;
		} else if (index===crownTitleIndex) {
			crownTitle = val;
		} else if (index===yearBornIndex) {
			yearBorn = val;
		} else if (index===yearDiedIndex) {
			yearDied = val;
		} else if (index===houseHeadIndex) {
			houseHead = val;
		} else if (index===reignStartIndex) {
			reignStart = val;
		} else if (index===reignEndIndex) {
			reignEnd = val;
		}
	});
	
	if (!memberName) {
		logger.logInfo(
				'No house member to update.\nArguments example: %s',
				'node UpdateGraph -memberName "Napoleon Bonaparte" -crownName "Grand Napoleon" -crownTitle "King of the world" -reignStart "8 Nov 2001" -yearDied 2016');
		return;
	}
		
	var conn = conFactory.getFactory().createConnection(url, user, pwd, null);
	conn.connect(function(){
		var gof = conn.getGraphObjectFactory();
		if (gof === null) {
			throw new TGException("Graph object not found");
		}

		conn.getGraphMetadata(true, function(){
			var houseKey = gof.createCompositeKey("houseMemberType");
			
			houseKey.setAttribute("memberName", memberName);
			logger.logInfo("Searching for member '%s'...",memberName);
      		var houseMember = conn.getEntity(houseKey, null, function(houseMember){
      			try {
              		if (houseMember) {
              			logger.logInfo("House member '%s' found",houseMember.getAttribute("memberName").getValue());
              			if (crownName) {
              				houseMember.setAttribute("crownName", crownName);
              			}
              			if (crownTitle) {
              				houseMember.setAttribute("crownTitle", crownTitle);
              			}
              			
              			if (houseHead) {
              			var isHouseHead = false;
              				if(houseHead==='true') {
              					isHouseHead = true;
              				} else if(houseHead==='false') {
              					isHouseHead = false;
              				} else if(houseHead==='null') {
              					isHouseHead = null;
              				} else {
              					throw new TGException('Illegel boolean value : houseHead -> ' + houseHead);
              				}
              				houseMember.setAttribute("houseHead", isHouseHead);
              			}
              			if (yearBorn) {
              				houseMember.setAttribute("yearBorn", parseInt(yearBorn));
              			}
              			if (yearDied) { 
              				if (yearDied === "null") {
              					houseMember.setAttribute("yearDied", null);
              				}
              				else {
              					houseMember.setAttribute("yearDied", parseInt(yearDied));
              				}
              			}

              			if (reignStart) {
              				if (reignStart==="null") {
              					houseMember.setAttribute("reignStart", null);
              				}
              				else {
              					try {
              						houseMember.setAttribute("reignStart", new Date(reignStart)); 
              					}
              					catch (exception) {
              						throw new TGException("Member update failed - Wrong parameter: -reignStart format should be \"dd MMM yyyy\"");
              					}
              				}
              			}
              			if (reignEnd) {
              				if (reignEnd==="null") {
              					houseMember.setAttribute("reignEnd", null);          					
              				} else {
              					try {
              						houseMember.setAttribute("reignEnd", new Date(reignEnd)); 
              					}
              					catch (exception) {
              						throw new TGException("Member update failed - Wrong parameter: -reignEnd format should be \"dd MMM yyyy\"");
              					}
              				}
              			}

              			conn.updateEntity(houseMember);
              			conn.commit(function(){
              				logger.logInfo("House member '%s' updated successfully", memberName);
                			if (conn) {
                				conn.disconnect();
                			}
              			});
              		} else {
              			logger.logInfo("House member '%s' not found", memberName);
              		}
      			} catch (exception) {
          			logger.logError("Exception happens, message ", exception.message);      				
         			if (conn) {
        				conn.disconnect();
        			}
      			}
     		});
		});
	});
}

main();
