#!/bin/bash

. scripts/utility.sh

function setEnvironmentVariables() {
  if [ $# -eq 0 ]; then
    Error "No arguments are supplied."
    return 1
  fi
  
  Info "Using organization config $1"
  if [ $1 == "Wales" ]; then
    export CORE_PEER_TLS_ENABLED=true
    export ORDERER_CA=${PWD}/crypto-config/ordererOrganizations/HA.com/orderers/orderer.HA.com/msp/tlscacerts/tlsca.HA.com-cert.pem
    export PEER0_ORG1_CA=${PWD}/crypto-config/peerOrganizations/WalesHospital.com/peers/peer0.WalesHospital.com/tls/ca.crt
    export PEER0_ORG2_CA=${PWD}/crypto-config/peerOrganizations/MargaretHospital.com/peers/peer0.MargaretHospital.com/tls/ca.crt
    export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/crypto-config/ordererOrganizations/HA.com/orderers/orderer.HA.com/tls/server.crt
    export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/crypto-config/ordererOrganizations/HA.com/orderers/orderer.HA.com/tls/server.key
    export CORE_PEER_LOCALMSPID="WalesHospitalMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/WalesHospital.com/users/Admin@WalesHospital.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    Info "Finished setting Environment Variable for Organization $1"
  elif [ $1 == "Margaret" ]; then
    export CORE_PEER_TLS_ENABLED=true
    export ORDERER_CA=${PWD}/crypto-config/ordererOrganizations/HA.com/orderers/orderer.HA.com/msp/tlscacerts/tlsca.HA.com-cert.pem
    export PEER0_ORG1_CA=${PWD}/crypto-config/peerOrganizations/WalesHospital.com/peers/peer0.WalesHospital.com/tls/ca.crt
    export PEER0_ORG2_CA=${PWD}/crypto-config/peerOrganizations/MargaretHospital.com/peers/peer0.MargaretHospital.com/tls/ca.crt
    export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/crypto-config/ordererOrganizations/HA.com/orderers/orderer.HA.com/tls/server.crt
    export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/crypto-config/ordererOrganizations/HA.com/orderers/orderer.HA.com/tls/server.key
    export CORE_PEER_LOCALMSPID="MargaretHospitalMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/MargaretHospital.com/users/Admin@MargaretHospital.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
    Info "Finished setting Environment Variable for Organization $1"
  else
    Error "Unknow Organization"
    return 1
  fi
}

