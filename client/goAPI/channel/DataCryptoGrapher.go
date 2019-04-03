package channel

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
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
 * File name: DataCryptoGrapher.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type DataCryptoGrapher struct {
	sessionId      int64
	remoteCert     *x509.Certificate
	pubKey         crypto.PublicKey
	algoParameters *pkix.AlgorithmIdentifier
}

func DefaultDataCryptoGrapher() *DataCryptoGrapher {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(DataCryptoGrapher{})

	newChannelUrl := DataCryptoGrapher{
		sessionId: 0,
	}

	return &newChannelUrl
}

func NewDataCryptoGrapher(sessionId int64, certbuffer []byte) (*DataCryptoGrapher, types.TGError) {
	newCryptoGrapher := DefaultDataCryptoGrapher()
	newCryptoGrapher.sessionId = sessionId

	// TODO: Uncomment once DataCryptoGrapher is implemented
	/**
	resultBlock, remainder := pem.Decode(certbuffer)
	if resultBlock.Type == "CERTIFICATE" {
		cert, err := x509.ParseCertificate(resultBlock.Bytes)
		if err != nil {
			errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse CERTIFICATE from the certificate buffer")
			return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		newCryptoGrapher.remoteCert = cert
	} else if resultBlock.Type == "PUBLIC KEY" {
		pubKey, err := x509.ParsePKIXPublicKey(resultBlock.Bytes)
		if err != nil {
			errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse PUBLIC KEY from the certificate buffer")
			return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		newCryptoGrapher.pubKey = pubKey
	}

	resultBlock, remainder = pem.Decode(remainder)
	if resultBlock.Type == "PUBLIC KEY" {
		pubKey, err := x509.ParsePKIXPublicKey(resultBlock.Bytes)
		if err != nil {
			errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse PUBLIC KEY from the certificate buffer")
			return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		newCryptoGrapher.pubKey = pubKey
	} else if resultBlock.Type == "CERTIFICATE" {
		cert, err := x509.ParseCertificate(resultBlock.Bytes)
		if err != nil {
			errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse CERTIFICATE from the certificate buffer")
			return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		newCryptoGrapher.remoteCert = cert
	}
	algoParams, err1 := getAlgorithmParameters(newCryptoGrapher.pubKey)
	if err1 != nil {
		errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse CERTIFICATE from the certificate buffer")
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err1.Error())
	}
	newCryptoGrapher.algoParameters = algoParams
	*/
	return newCryptoGrapher, nil
}

/////////////////////////////////////////////////////////////////
// Helper functions for DataCryptoGrapher
/////////////////////////////////////////////////////////////////

func (obj *DataCryptoGrapher) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DataCryptoGrapher:{")
	buffer.WriteString(fmt.Sprintf("SessionId: %d", obj.sessionId))
	//buffer.WriteString(fmt.Sprintf(", RemoteCert: %+v", obj.remoteCert))
	//buffer.WriteString(fmt.Sprintf(", PubKey: %+v", obj.pubKey))
	//buffer.WriteString(fmt.Sprintf(", AlgoParameters: %d", obj.algoParameters))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Private functions for DataCryptoGrapher
/////////////////////////////////////////////////////////////////

func getAlgorithmParameters(pubKey crypto.PublicKey) (*pkix.AlgorithmIdentifier, types.TGError) {
	// TODO: Uncomment once DataCryptoGrapher is implemented
	/**
			if (publicKey == null) return null;

	        if (publicKey instanceof DSAPublicKey) {
	            AlgorithmParameters algparams = AlgorithmParameters.getInstance(publicKey.getAlgorithm());
	            DSAPublicKey dsakey = (DSAPublicKey) publicKey;
	            DSAParameterSpec dsaParams = (DSAParameterSpec) dsakey.getParams();
	            algparams.init(dsaParams);
	            return algparams;
	        }

	        if (publicKey instanceof ECPublicKey) {
	            AlgorithmParameters algparams = AlgorithmParameters.getInstance(publicKey.getAlgorithm());
	            ECPublicKey eckey = (ECPublicKey) publicKey;
	            ECParameterSpec ecParams = (ECParameterSpec) eckey.getParams();
	            algparams.init(ecParams);
	            return algparams;
	        }

	        if (publicKey instanceof DHPublicKey) {
	            AlgorithmParameters algparams = AlgorithmParameters.getInstance(publicKey.getAlgorithm());
	            DHPublicKey dhkey = (DHPublicKey) publicKey;
	            DHParameterSpec dhParams = (DHParameterSpec) dhkey.getParams();
	            algparams.init(dhParams);
	            return algparams;
	        }

	        return null;  //RSA doesn't have
	*/
	// No-op for Now!
	return nil, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGDataCryptoGrapher
/////////////////////////////////////////////////////////////////

// Decrypt decrypts the buffer
func (obj *DataCryptoGrapher) Decrypt(encBuffer []byte) ([]byte, types.TGError) {
	// No-op for Now!
	return make([]byte, 0), nil
}

// Encrypt encrypts the buffer
func (obj *DataCryptoGrapher) Encrypt(decBuffer []byte) ([]byte, types.TGError) {
	// TODO: Uncomment once DataCryptoGrapher is implemented
	/**
		try {
			Cipher cipher = Cipher.getInstance(publicKey.getAlgorithm());
			cipher.init(Cipher.ENCRYPT_MODE, publicKey, algparams);
			return cipher.doFinal(data);
		}
		catch (Exception e) {
			throw new TGException(e);
		}
	*/
	return nil, nil
}
