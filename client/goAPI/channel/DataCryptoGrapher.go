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

package channel

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"golang.org/x/crypto/blowfish"
)

type DataCryptoGrapher struct {
	sessionId      int64
	remoteCert     *x509.Certificate
	pubKey         crypto.PublicKey
	//algoParameters *pkix.AlgorithmIdentifier
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

func NewDataCryptoGrapher(sessionId int64, serverCertBytes []byte) (*DataCryptoGrapher, types.TGError) {
	logger.Log(fmt.Sprintf("Entering NewDataCryptoGrapher() w/ serverCertBytes as '%+v'", serverCertBytes))
	newCryptoGrapher := DefaultDataCryptoGrapher()
	newCryptoGrapher.sessionId = sessionId

	// TODO: Uncomment once DataCryptoGrapher is implemented
	//logger.Log(fmt.Sprintf("Inside NewDataCryptoGrapher - about to x509.ParseCertificate(()"))
	//cert, err := x509.ParseCertificate(serverCertBytes)
	//if err != nil {
	//	errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse CERTIFICATE from the certificate buffer")
	//	return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
	//}
	//newCryptoGrapher.remoteCert = cert
	//logger.Log(fmt.Sprintf("Inside NewDataCryptoGrapher - parsed certificate as'%+v'", cert))
	//
	//logger.Log(fmt.Sprint("Inside NewDataCryptoGrapher - about to x509.ParsePKIXPublicKey()"))
	//pubKey, err := x509.ParsePKIXPublicKey(serverCertBytes)
	//if err != nil {
	//	errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse PUBLIC KEY from the certificate buffer")
	//	return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
	//}
	//newCryptoGrapher.pubKey = pubKey
	//logger.Log(fmt.Sprintf("Inside NewDataCryptoGrapher - parsed public key as'%+v'", pubKey))

	/**
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
	buffer.WriteString(fmt.Sprintf(", RemoteCert: %+v", obj.remoteCert))
	buffer.WriteString(fmt.Sprintf(", PubKey: %+v", obj.pubKey))
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
//func (obj *DataCryptoGrapher) Decrypt(encBuffer []byte) ([]byte, types.TGError) {
func (obj *DataCryptoGrapher) Decrypt(is types.TGInputStream) ([]byte, types.TGError) {
	logger.Log(fmt.Sprint("Entering DataCryptoGrapher:Decrypt()"))
	out := iostream.DefaultProtocolDataOutputStream()
	buf := make([]byte, 0)

	rand, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DataCryptoGrapher:Decrypt w/ Error in reading rand from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside DataCryptoGrapher:Decrypt read resultId as '%+v'", rand))

	len, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DataCryptoGrapher:Decrypt w/ Error in reading len from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside DataCryptoGrapher:Decrypt read resultId as '%+v'", len))

	cnt := len / 8
	rem := len % 8

	for i:=0; i<int(cnt); i++ {
		val, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning DataCryptoGrapher:Decrypt w/ Error in reading val from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside DataCryptoGrapher:Decrypt read resultId as '%+v'", val))

		org := val ^ rand
		_ = out.WriteLongAsBytes(org)
	}

	for i:=0; i<int(rem); i++ {
		val, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning DataCryptoGrapher:Decrypt w/ Error in reading val from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside DataCryptoGrapher:Decrypt read resultId as '%+v'", val))

		out.WriteByte(int(val))
	}
	logger.Log(fmt.Sprintf("Returning DataCryptoGrapher:Decrypt() w/ decrypted buffer as '%+v'", buf))
	return out.ToByteArray()
}

// Encrypt encrypts the buffer
func (obj *DataCryptoGrapher) Encrypt(rawBuf []byte) ([]byte, types.TGError) {
	logger.Log(fmt.Sprintf("Entering DataCryptoGrapher:Encrypt() w/ raw buffer as '%+v'", rawBuf))
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

	/**
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs5Padding()
	pt, err = padder.Pad(pt) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
	*/

	block, err := blowfish.NewCipher(obj.remoteCert.RawSubjectPublicKeyInfo)
	if err != nil {
		return nil, exception.GetErrorByType(types.TGErrorSecurityException, types.INTERNAL_SERVER_ERROR, err.Error(), "")
	}
	encryptedBuf := make([]byte, aes.BlockSize+len(rawBuf))
	iv := encryptedBuf[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptedBuf[aes.BlockSize:], rawBuf)
	fmt.Printf("%x\n", encryptedBuf)

	/**
	block := blowfish.NewCipher(obj.remoteCert.RawSubjectPublicKeyInfo)
	encryptedBuf := make([]byte, aes.BlockSize+len(decBuffer))
	iv := encryptedBuf[:aes.BlockSize]

	algo := obj.remoteCert.PublicKeyAlgorithm
	switch algo {
	case x509.RSA:
	case x509.DSA:
		block, err := des.NewCipher(rawBuf)
		if err != nil {
			return nil, exception.GetErrorByType(types.TGErrorSecurityException, types.INTERNAL_SERVER_ERROR, err.Error(), "")
		}
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(encryptedBuf[aes.BlockSize:], rawBuf)
	case x509.ECDSA:
	}
	*/
	//logger.Log(fmt.Sprintf("Returning DataCryptoGrapher:Decrypt() w/ encrypted buffer as '%+v'", encryptedBuf))
	return nil, nil
}
