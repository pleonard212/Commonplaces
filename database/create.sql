CREATE TABLE `alignments` (
  `sourcephiloid` varchar(255) DEFAULT NULL, 				        /* primary key of doc 1 */
  `targetphiloid` varchar(255) DEFAULT NULL,				        /* primary key of doc 2 */
  `sourceleftcontext` text,								        	        /* text to the left of match in doc 1 */
  `sourcematchcontext` text,							        	        /* match in doc 1 */
  `sourcerightcontext` text,						        		        /* text to the right of match in doc 1 */
  `targetleftcontext` text,							        		        /* text to the left of match in doc 2 */
  `targetmatchcontext` text,						        		        /* match in doc 2 */
  `targetrightcontext` text,							        	        /* text to the right of match in doc 2 */
  `sourceauthor` varchar(128) DEFAULT NULL,					        /* doc 1 author */
  `sourcetitle` varchar(256) DEFAULT NULL,					        /* doc 1 title */
  `sourcedate` int(4) DEFAULT NULL,							            /* doc 1 year */
  `sourcemodulename` varchar(128) DEFAULT NULL,			        /* flexible; think of this as doc 1 genre (novel, poem) or ECCO subject category */
  `targetauthor` varchar(128) DEFAULT NULL,					        /* doc 2 author */			
  `targettitle` varchar(256) DEFAULT NULL,				        	/* doc 2 title */
  `targetdate` int(4) DEFAULT NULL,							            /* doc 1 year */
  `targetmodulename` varchar(128) DEFAULT NULL,			        /* flexible; think of this as doc 2 genre (novel, poem) or ECCO subject category */
  `authorident` int(32) DEFAULT NULL,					     	        /* this gets used to flag books as the Bible or "Book of Common Expressions" etc */
  `passageident` int(64) unsigned NOT NULL AUTO_INCREMENT,	/* the unique id of the match */
  `passageidentcount` int(32) DEFAULT NULL,					        /* Hmm, I dunno. Have to inspect the Go code. */
  PRIMARY KEY (`passageident`)
) ENGINE=MyISAM AUTO_INCREMENT=68 DEFAULT CHARSET=utf8;     /* MyISAM is actually important for the search strategy. */
