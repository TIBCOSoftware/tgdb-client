package com.tibco.tgdb.exception;

/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except 
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name :TGBadMagic
 * Created by: suresh
 *
 * SVN Id: $Id: TGBadMagic.java 622 2016-03-19 20:51:12Z ssubrama $
 */
public class TGBadMagic extends TGException {

    public TGBadMagic(String reason, String err) {
        super(reason, err);
    }

    @Override
    public ExceptionType getExceptionType() {
        return ExceptionType.BadMagic;
    }
}
