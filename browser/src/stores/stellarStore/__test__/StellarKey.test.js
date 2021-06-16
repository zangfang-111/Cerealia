import stellarStore from '../index'

let keypairs = {
  'k1': ['SCRHEWSED4VXPCA55HBNWSZ7ESRU2NLLNJRIO532DJBPSB5X3F7UVNRM', 'GDOKCE5VFBB3CCPWG6HLQXW7AL4QELDSHHLABWDBYMRSI4UGYY4BBGS3'],
  'k2': ['SDOM5RGWXOYFKXLOL4OKUQ7EAKO4KVCTLEZCQCESIAVQOQFMPRWOKAKY', 'GCZVKTLPQY54OGK5R3EEAU22Q7COES2XP5544H2C3CRNN7GGRL74IBUM'],
  'k3': ['SCVZ2D7TXSH7RNIYJOQSJOBTJSWE7AMDOHSVBQJA3OZABTJQ5MC47XID', 'GD3EPS4EBOK6ZELDEN466I6EU4LW7TK6UL6INZRC6OKLZKYXXESS64VE']
}

function validateKey (keypair, expected) {
  expect(stellarStore.validateStellarSecretKey(keypair[0])).toBe(expected)
  expect(stellarStore.validateStellarPublicKey(keypair[1])).toBe(expected)
  expect(stellarStore.validateAndSetUserKey(keypair[0], keypair[1])).toBe(expected)
}

test('test for keyPair sample 1', () => {
  validateKey(keypairs.k1, true)
})

it('test for keyPair sample 2', () => {
  validateKey(keypairs.k2, true)
})

it('test for keyPair sample 3', () => {
  validateKey(keypairs.k3, true)
})

it('test for wrong keyPair', () => {
  let pk = keypairs.k1[1]
  let sk = keypairs.k1[0]
  sk = '0' + sk.substr(1)
  pk = '0' + pk.substr(1)
  validateKey([sk, pk], false)
})
