#ifndef _FCITX5_BINARYIME_BINARYIME_H_
#define _FCITX5_BINARYIME_BINARYIME_H_

#include <fcitx/inputmethodengine.h>
#include <fcitx/addonfactory.h>
#include <fcitx/inputcontextproperty.h>
#include <fcitx/inputpanel.h>
#include <fcitx/event.h>
#include <fcitx/instance.h>
#include <fcitx/addonmanager.h>
#include <unordered_map>
#include <string>

class BinaryIMEState : public fcitx::InputContextProperty {
public:
    std::string currentText;
    int dualMode = 0; // 0=BINARY, 1=TOGGLE, 2=PRINT
    bool ignoreNextKeyword = false;
    std::string pendingConversion;
    
    void copyTo(fcitx::InputContextProperty* other) override {
        auto* otherState = static_cast<BinaryIMEState*>(other);
        otherState->dualMode = dualMode;
        otherState->ignoreNextKeyword = ignoreNextKeyword;
    }
    
    bool needCopy() const override { return true; }
};

class BinaryIME : public fcitx::InputMethodEngineV2 {
public:
    BinaryIME(fcitx::Instance* instance);
    
    void keyEvent(const fcitx::InputMethodEntry& entry, 
                  fcitx::KeyEvent& keyEvent) override;
    
private:
    void loadConfig();
    bool loadConfigFromFile(const std::string& filename);
    void loadDefaultConfig();
    
    std::string convertToBinary(const std::string& text, fcitx::InputContext* ic = nullptr);
    std::string convertTextToBinary(const std::string& text);
    std::string convertSentenceToBinary(const std::string& sentence);
    std::string convertSentenceInPrintMode(const std::string& text, bool ignoreKeyword);
    std::string convertNumber(const std::string& numberStr);
    std::string processWord(const std::string& word);
    
    bool isNumber(const std::string& str);
    void updatePreedit(fcitx::InputContext* ic);
    void commitConversion(fcitx::InputContext* ic);
    
    std::unordered_map<std::string, std::string> charMap;
    std::unordered_map<std::string, std::string> capitalMap; 
    std::unordered_map<std::string, std::string> digitsMap;
    std::unordered_map<std::string, std::string> keywords;
    std::unordered_map<std::string, std::string> operators;
    std::unordered_map<std::string, std::string> specialMap; 


    bool convertNumbers = true;
    bool addSpaces = true;
    bool caseSensitive = false; 
    std::string unknownBehavior = "keep";
    std::string numberSeparator = "$";


    fcitx::FactoryFor<BinaryIMEState> stateFactory_;
    fcitx::Instance* instance_;
};

class BinaryIMEFactory : public fcitx::AddonFactory {
public:
    fcitx::AddonInstance* create(fcitx::AddonManager* manager) override {
        return new BinaryIME(manager->instance());
    }
};

#endif // _FCITX5_BINARYIME_BINARYIME_H_