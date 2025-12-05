#include "{{.IMEName}}.h"
#include <fcitx/inputcontext.h>
#include <fcitx-utils/key.h>
#include <fstream>
#include <iostream>
#include <cctype>
#include <algorithm>
#include <sstream> 

{{.IMEName}}::{{.IMEName}}(fcitx::Instance* instance) 
    : stateFactory_([](fcitx::InputContext&) { return new {{.IMEName}}State; }),
      instance_(instance) {
    
    instance_->inputContextManager().registerProperty("{{.IMEName}}State", &stateFactory_);
    loadConfig();
    FCITX_INFO() << "{{.IMEName}} initialized successfully!";
}

void {{.IMEName}}::loadConfig() {
    FCITX_INFO() << "Loading {{.IMEName}} configuration...";
    
    std::string homeDir = std::getenv("HOME");
    std::string userConfigFile = homeDir + "/.config/custom-ime/{{.ProjectName}}/config/{{.IMEName}}.conf";
    
    std::string systemConfigFile = "/usr/share/fcitx5/{{.ProjectName}}/config/{{.IMEName}}.conf";
    
    bool configLoaded = false;
    
    if (loadConfigFromFile(userConfigFile)) {
        FCITX_INFO() << "Loaded user config from: " << userConfigFile;
        configLoaded = true;
    } 
    else if (loadConfigFromFile(systemConfigFile)) {
        FCITX_INFO() << "Loaded system config from: " << systemConfigFile;
        configLoaded = true;
    }
    
    if (!configLoaded) {
        FCITX_WARN() << "No config file found, using default configuration";
        FCITX_WARN() << "Searched in:";
        FCITX_WARN() << "  " << userConfigFile;
        FCITX_WARN() << "  " << systemConfigFile;
        loadDefaultConfig();
    }
    
    FCITX_INFO() << "Configuration loaded - Chars: " << charMap.size() 
                 << ", Capitals: " << capitalMap.size()
                 << ", Keywords: " << keywords.size() 
                 << ", Operators: " << operators.size()
                 << ", Digits: "<< digitsMap.size()
                 << ", Convert Numbers: " << (convertNumbers ? "true" : "false")
                 << ", Separator: '" << numberSeparator << "'";
}

bool {{.IMEName}}::loadConfigFromFile(const std::string& filename) {
    std::ifstream file(filename);
    if (!file.is_open()) {
        FCITX_ERROR() << "Cannot open config file: " << filename;
        return false;
    }
    
    std::string currentSection;
    std::string line;
    
    while (std::getline(file, line)) {
        size_t commentPos = line.find('#');
        if (commentPos != std::string::npos) {
            line = line.substr(0, commentPos);
        }
        
        line.erase(0, line.find_first_not_of(" \t"));
        line.erase(line.find_last_not_of(" \t") + 1);
        
        if (line.empty()) continue;
        
        if (line[0] == '[' && line[line.length()-1] == ']') {
            currentSection = line.substr(1, line.length()-2);
            FCITX_INFO() << "Loading section: " << currentSection;
            continue;
        }
        
        size_t equalPos = line.find('=');
        if (equalPos == std::string::npos) continue;
        
        std::string key = line.substr(0, equalPos);
        std::string value = line.substr(equalPos + 1);
        
        key.erase(0, key.find_first_not_of(" \t"));
        key.erase(key.find_last_not_of(" \t") + 1);
        value.erase(0, value.find_first_not_of(" \t"));
        value.erase(value.find_last_not_of(" \t") + 1);
        
        if (key.empty() || value.empty()) continue;
        
        if (currentSection == "Settings") {
            if (key == "convert_numbers_to_binary") {
                convertNumbers = (value == "true");
            } else if (key == "unknown_chars_behavior") {
                unknownBehavior = value;
            } else if (key == "add_spaces") {
                addSpaces = (value == "true");
            } else if (key == "number_separator") {
                numberSeparator = value;
            } else if (key == "case_sensitive") {
                caseSensitive = (value == "true");
            }
        }
        else if (currentSection == "Characters") {
            charMap[key] = value;
        }
        else if (currentSection == "Capitals") {
            capitalMap[key] = value;
        }
        else if (currentSection == "Digits") {
            FCITX_INFO() << "got digit key : " << key << " value : " <<value ;
            digitsMap[key] = value;
        }
        else if (currentSection == "Keywords") {
            keywords[key] = value;
        }
        else if (currentSection == "Operators") {
            operators[key] = value;
        }
        else if (currentSection == "Special") {
            specialMap[key] = value;
        }
    }
    
    file.close();
    return true;
}

void {{.IMEName}}::loadDefaultConfig() {
    convertNumbers = true;
    unknownBehavior = "keep";
    addSpaces = true;
    numberSeparator = "$";
    caseSensitive = false;
    
    charMap = {
        {"a", "00000"}, {"b", "00001"}, {"c", "00010"}, {"d", "00011"},
        {"e", "00100"}, {"f", "00101"}, {"g", "00110"}, {"h", "00111"},
        {"i", "01000"}, {"j", "01001"}, {"k", "01010"}, {"l", "01011"},
        {"m", "01100"}, {"n", "01101"}, {"o", "01110"}, {"p", "01111"},
        {"q", "10000"}, {"r", "10001"}, {"s", "10010"}, {"t", "10011"},
        {"u", "10100"}, {"v", "10101"}, {"w", "10110"}, {"x", "10111"},
        {"y", "11000"}, {"z", "11001"}
    };
    
    capitalMap = {
        {"A", "00000"}, {"B", "00001"}, {"C", "00010"}, {"D", "00011"},
        {"E", "00100"}, {"F", "00101"}, {"G", "00110"}, {"H", "00111"},
        {"I", "01000"}, {"J", "01001"}, {"K", "01010"}, {"L", "01011"},
        {"M", "01100"}, {"N", "01101"}, {"O", "01110"}, {"P", "01111"},
        {"Q", "10000"}, {"R", "10001"}, {"S", "10010"}, {"T", "10011"},
        {"U", "10100"}, {"V", "10101"}, {"W", "10110"}, {"X", "10111"},
        {"Y", "11000"}, {"Z", "11001"}
    };
    
    digitsMap = {
        {"0", "0"}, {"1", "1"}, {"2", "10"}, {"3", "11"}, {"4", "100"}, 
        {"5", "101"}, {"6", "110"}, {"7", "111"}, {"8", "1000"}, {"9", "1001"}
    };
    
    keywords = {
        {"if", "001"}, {"else", "010"}, {"for", "011"}, {"while", "100"},
        {"def", "101"}, {"return", "110"}, {"class", "111"}
    };
    
    operators = {
        {"==", "00"}, {"!=", "01"}, {"=", "10"}, {"+", "11"}
    };
}

void {{.IMEName}}::keyEvent(const fcitx::InputMethodEntry& entry, 
                         fcitx::KeyEvent& keyEvent) {
    FCITX_UNUSED(entry);
    
    auto* ic = keyEvent.inputContext();
    if (!ic) return;
    
    auto* state = ic->propertyFor(&stateFactory_);
    
    if (keyEvent.isRelease()) {
        return;
    }
    
    auto key = keyEvent.key();
    
    if (key.check(FcitxKey_F8)) {
        keyEvent.filterAndAccept();
        state->dualMode = (state->dualMode + 1) % 3;
        std::string modeNames[] = {"ORG", "TOGGLE", "PRINT"};
        FCITX_INFO() << "Mode: " << modeNames[state->dualMode];
        updatePreedit(ic);
        return;
    }
    
    if (key.check(FcitxKey_F7)) {
        keyEvent.filterAndAccept();
        state->ignoreNextKeyword = !state->ignoreNextKeyword;
        FCITX_INFO() << "Ignore keyword: " << (state->ignoreNextKeyword ? "ON" : "OFF");
        updatePreedit(ic);
        return;
    }
    
    if (key.check(FcitxKey_space)) {
        keyEvent.filterAndAccept();
        
        if (state->dualMode == 2) { // PRINT Mode
            state->currentText += " ";
            updatePreedit(ic);
        } else { 
            if (!state->currentText.empty()) {
                std::string converted = convertSentenceToBinary(state->currentText);
                ic->commitString(converted + " ");
                FCITX_INFO() << "Converted: '" << state->currentText << "' -> '" << converted << "'";
                state->currentText.clear();
                updatePreedit(ic);
            } else {
                ic->commitString(" ");
            }
        }
        return;
    }
    
    if (key.check(FcitxKey_BackSpace)) {
        if (!state->currentText.empty()) {
            keyEvent.filterAndAccept();
            state->currentText.pop_back();
            updatePreedit(ic);
            FCITX_INFO() << "Backspace in preedit - Current text: '" << state->currentText << "'";
        } else {
            keyEvent.filter();
            FCITX_INFO() << "Backspace passed to application";
        }
        return;
    }
   

    if (key.check(FcitxKey_Return) || key.check(FcitxKey_KP_Enter)) {
        if (!state->currentText.empty()) {
            keyEvent.filterAndAccept();
            std::string converted;
            std::string finalOutput;
            
            if (state->dualMode == 2) { // PRINT Mode
                converted = convertToBinary(state->currentText, ic);
                finalOutput = state->currentText + "\n" + converted;
            } else { 
                converted = convertSentenceToBinary(state->currentText);
                finalOutput = converted;
            }
            
            ic->commitString(finalOutput);
            FCITX_INFO() << "Converted and committed: '" << state->currentText << "' -> '" << converted << "'";
            state->currentText.clear();
            updatePreedit(ic);
        } else {
            keyEvent.filter();
            return;
        }
        return;
    }
    
    if (key.isSimple()) {
        keyEvent.filterAndAccept();
        uint32_t sym = key.sym();
        
        if ((sym >= 'a' && sym <= 'z') || 
            (sym >= 'A' && sym <= 'Z') ||
            (sym >= '0' && sym <= '9') ||
            sym == '$' || std::ispunct(sym)) {
            
            char c = static_cast<char>(sym);
            state->currentText += c;
            updatePreedit(ic);
            FCITX_INFO() << "Current text: " << state->currentText;
        }
        return;
    }
    
    keyEvent.filter();
}

std::string {{.IMEName}}::convertToBinary(const std::string& text, fcitx::InputContext* ic) {
    if (text.empty()) return "";
    
    FCITX_INFO() << "Converting: '" << text << "' (ConvertNumbers: " << convertNumbers << ")";
    
    bool ignoreKeyword = false;
    if (ic) {
        auto* state = ic->propertyFor(&stateFactory_);
        ignoreKeyword = state->ignoreNextKeyword;
    }
    
    if (ic) {
        auto* state = ic->propertyFor(&stateFactory_);
        if (state->dualMode == 2) { //  PRINT Mode
            if (!ignoreKeyword && keywords.find(text) != keywords.end()) {
                FCITX_INFO() << "Found keyword: " << text << " -> " << keywords[text];
                return keywords[text];
            }
            
            if (operators.find(text) != operators.end()) {
                FCITX_INFO() << "Found operator: " << text << " -> " << operators[text];
                return operators[text];
            }
            
            if (isNumber(text)) {
                return convertNumber(text);
            }
            
            if (!text.empty() && text[0] == numberSeparator[0] && text.size() > 1) {
                std::string numberPart = text.substr(1);
                if (isNumber(numberPart)) {
                    return convertNumber(numberPart);
                }
            }
            
            if (text.find(' ') != std::string::npos) {
                return convertSentenceInPrintMode(text, ignoreKeyword);
            }
        }
    }
    
    return convertTextToBinary(text);
}

std::string {{.IMEName}}::convertNumber(const std::string& numberStr) {
    if (!convertNumbers) {
        std::string result;
        
        std::string cleanNumber = numberStr;
        if (!cleanNumber.empty() && cleanNumber[0] == numberSeparator[0]) {
            cleanNumber = cleanNumber.substr(1);
        }
        
        if (isNumber(cleanNumber)) {
            if (specialMap.find("number_sign") != specialMap.end()) {
                result += specialMap["number_sign"];
            } else {
                result += numberSeparator;
            }
            
            for (char c : cleanNumber) {
                std::string digit(1, c);
                if (digitsMap.find(digit) != digitsMap.end()) {
                    result += digitsMap[digit];
                } else {
                    result += digit;
                }
            }
            FCITX_INFO() << "Converted number '" << numberStr << "' -> '" << result << "'";
            return result;
        } else {
            return numberStr;
        }
    } else {
        int number = std::stoi(numberStr);
        std::string binary;
        if (number == 0) {
            binary = "0";
        } else {
            while (number > 0) {
                binary = (number % 2 ? "1" : "0") + binary;
                number /= 2;
            }
        }
        return numberSeparator + binary;
    }
}

std::string {{.IMEName}}::convertSentenceInPrintMode(const std::string& text, bool ignoreKeyword) {
    std::stringstream ss(text);
    std::string word;
    std::string result;
    
    while (ss >> word) {
        if (!result.empty()) {
            result += " ";
        }
        
        std::string convertedWord;
        
        if (!ignoreKeyword && keywords.find(word) != keywords.end()) {
            convertedWord = keywords[word];
        }
        else if (operators.find(word) != operators.end()) {
            convertedWord = operators[word];
        }
        else if (isNumber(word) ) {
            FCITX_INFO() << "here1 convert numebr ";
            convertedWord = convertNumber(word);
        }
        else if (!word.empty() && word[0] == numberSeparator[0] && word.size() > 1 )  {
            FCITX_INFO() << "here12convert numebr ";

            std::string numberPart = word.substr(1);
            if (isNumber(numberPart)) {
                convertedWord = convertNumber(numberPart);
            } else {
                convertedWord = convertTextToBinary(word);
            }
        }
        else {
            FCITX_INFO() << "here3 convert numebr ";

            convertedWord = convertTextToBinary(word);
        }
        
        result += convertedWord;
    }
    
    FCITX_INFO() << "PRINT Mode sentence conversion: '" << text << "' -> '" << result << "'";
    return result;
}

std::string {{.IMEName}}::convertTextToBinary(const std::string& text) {
    if (text.empty()) return "";
    
    FCITX_INFO() << "Converting text: '" << text << "' (CaseSensitive: " << caseSensitive << ")";
    
    std::string result;
    std::string currentNumber;
    
    for (size_t i = 0; i < text.length(); ++i) {
        char c = text[i];
        std::string charStr(1, c);
        std::string lowerCharStr(1, std::tolower(c));
        
        if (std::isdigit(c)) {
            currentNumber += c;
            if (i == text.length() - 1 || !std::isdigit(text[i + 1])) {
                if (!currentNumber.empty()) {
                    result += convertNumber(currentNumber);
                    currentNumber.clear();
                }
            }
        } 
        else if (caseSensitive && std::isupper(c) && capitalMap.find(charStr) != capitalMap.end()) {
            if (!currentNumber.empty()) {
                result += convertNumber(currentNumber);
                currentNumber.clear();
            }
            result += capitalMap[charStr];
        }
        else if (charMap.find(lowerCharStr) != charMap.end()) {
            if (!currentNumber.empty()) {
                result += convertNumber(currentNumber);
                currentNumber.clear();
            }
            result += charMap[lowerCharStr];
        }
        else if (charMap.find(charStr) != charMap.end()) {
            if (!currentNumber.empty()) {
                result += convertNumber(currentNumber);
                currentNumber.clear();
            }
            result += charMap[charStr];
        }
        else if (c == ' ') {
            if (!currentNumber.empty()) {
                result += convertNumber(currentNumber);
                currentNumber.clear();
            }
            if (specialMap.find("space") != specialMap.end()) {
                if (!result.empty() && result.back() != ' ') {
                    result += specialMap["space"];
                }
            } else {
                if (!result.empty() && result.back() != ' ') {
                    result += " ";
                }
            }
        } 
        else {
            if (!currentNumber.empty()) {
                result += convertNumber(currentNumber);
                currentNumber.clear();
            }
            if (unknownBehavior == "keep"){ 
                result += charStr;
            }
            // ignore the input char if unknownBehavior != "keep"
        }
    }
    
    if (!currentNumber.empty()) {
        result += convertNumber(currentNumber);
    }
    
    FCITX_INFO() << "Text conversion: '" << text << "' -> '" << result << "'";
    return result;
}

std::string {{.IMEName}}::convertSentenceToBinary(const std::string& sentence) {
    if (sentence.empty()) return "";
    
    FCITX_INFO() << "Converting sentence: '" << sentence << "'";
    
    auto* ic = instance_->lastFocusedInputContext();
    if (ic) {
        auto* state = ic->propertyFor(&stateFactory_);
        if (!state->ignoreNextKeyword && keywords.find(sentence) != keywords.end()) {
            FCITX_INFO() << "Found keyword: " << sentence << " -> " << keywords[sentence];
            return keywords[sentence];
        }
    }
    
    if (operators.find(sentence) != operators.end()) {
        FCITX_INFO() << "Found operator: " << sentence << " -> " << operators[sentence];
        return operators[sentence];
    }
    
    if (isNumber(sentence)) {
        return convertNumber(sentence);
    }
    
    if (!sentence.empty() && sentence[0] == numberSeparator[0] && sentence.size() > 1) {
        std::string numberPart = sentence.substr(1);
        if (isNumber(numberPart)) {
            return convertNumber(numberPart);
        }
    }
    
    std::stringstream ss(sentence);
    std::string word;
    std::string result;
    
    while (ss >> word) {
        if (!result.empty()) {
            result += " ";
        }
        result += convertToBinary(word, nullptr);
    }
    
    FCITX_INFO() << "Sentence conversion: '" << sentence << "' -> '" << result << "'";
    return result;
}

std::string {{.IMEName}}::processWord(const std::string& word) {
    return convertToBinary(word, nullptr);
}

bool {{.IMEName}}::isNumber(const std::string& str) {
    return !str.empty() && std::all_of(str.begin(), str.end(), ::isdigit);
}

void {{.IMEName}}::updatePreedit(fcitx::InputContext* ic) {
    auto* state = ic->propertyFor(&stateFactory_);
    
    if (state->currentText.empty()) {
        ic->inputPanel().reset();
        ic->updatePreedit();
        ic->updateUserInterface(fcitx::UserInterfaceComponent::InputPanel);
        return;
    }
    
    std::string converted;
    if (state->dualMode == 2) { 
        converted = convertToBinary(state->currentText, ic);
    } else {
        converted = convertSentenceToBinary(state->currentText);
    }
    
    fcitx::Text preeditText;
    std::string modeNames[] = {"ORG", "TOGGLE", "PRINT"};
    
    if (state->dualMode == 1 || state->dualMode == 2) { 
        preeditText.append(state->currentText + "    [Mode:" + modeNames[state->dualMode] + 
                          " Ignore:" + (state->ignoreNextKeyword ? "ON" : "OFF") + "]");
        preeditText.append("\n" + converted);
    } else { 
        preeditText.append(converted + "    [Mode:" + modeNames[state->dualMode] + 
                          " Ignore:" + (state->ignoreNextKeyword ? "ON" : "OFF") + "]");
    }
    
    ic->inputPanel().setPreedit(preeditText);
    ic->updatePreedit();
    ic->updateUserInterface(fcitx::UserInterfaceComponent::InputPanel);
}

void {{.IMEName}}::commitConversion(fcitx::InputContext* ic) {
    auto* state = ic->propertyFor(&stateFactory_);
    if (!state->currentText.empty()) {
        std::string converted;
        std::string finalOutput;
        
        if (state->dualMode == 2) { // PRINT Mode
            converted = convertToBinary(state->currentText, ic);
            finalOutput = state->currentText + "\n" + converted;
        } else {
            converted = convertSentenceToBinary(state->currentText);
            finalOutput = converted;
        }
        
        ic->commitString(finalOutput);
        FCITX_INFO() << "Committed: '" << finalOutput << "'";
        state->currentText.clear();
        updatePreedit(ic);
    }
}

FCITX_ADDON_FACTORY({{.IMEName}}Factory)