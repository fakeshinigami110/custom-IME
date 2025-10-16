import re
from config_loader import load_mapping_config, load_engine_config, load_keywords_config


class BinaryConverter:
    def __init__(self):
        self.mapping = load_mapping_config()
        self.engine_config = load_engine_config()
        self.keywords, self.operators = load_keywords_config()

        print("Loaded operators:", self.operators)  # برای دیباگ

    def convert_text(self, text):
        """تبدیل متن به باینری"""
        if not text:
            return ""

        # ابتدا عملگرها را جایگزین کن
        for op, binary in sorted(self.operators.items(), key=lambda x: len(x[0]), reverse=True):
            text = text.replace(op, f" {binary} ")

        # سپس کلمات کلیدی را جایگزین کن
        words = text.split()
        result = []

        for word in words:
            # اگر از قبل باینری است (عملگر)
            if all(c in '01 ' for c in word):
                result.append(word)
            # کلمات کلیدی
            elif word in self.keywords:
                result.append(self.keywords[word])
            # اعداد
            elif self._is_number(word):
                result.append('$' + self._number_to_binary(word))
            # حروف و کلمات
            else:
                result.append(self._convert_letters(word))

        return ' '.join(result)

    def _is_number(self, word):
        """بررسی آیا کلمه عدد است"""
        if word.startswith('-'):
            rest = word[1:]
        else:
            rest = word

        if rest.replace('.', '').isdigit() and rest.count('.') <= 1:
            return True
        return False

    def _convert_letters(self, text):
        """تبدیل حروف به باینری"""
        result = ""
        for char in text:
            if char.lower() in self.mapping:
                result += self.mapping[char.lower()]
            else:
                result += char
        return result

    def _number_to_binary(self, number_str):
        """تبدیل عدد به باینری"""
        try:
            is_negative = number_str.startswith('-')
            clean_number = number_str[1:] if is_negative else number_str

            if '.' in clean_number:
                integer_part, fractional_part = clean_number.split('.')
                integer_bin = bin(int(integer_part))[2:] if integer_part else "0"
                fractional_bin = self._fraction_to_binary(fractional_part)
                result = f"{integer_bin}.{fractional_bin}"
            else:
                result = bin(int(clean_number))[2:]

            return f"0 {result}" if is_negative else result

        except (ValueError, Exception):
            return number_str

    def _fraction_to_binary(self, fractional_str, precision=8):
        """تبدیل بخش اعشاری به باینری"""
        try:
            if not fractional_str or fractional_str == "0":
                return "0"

            fractional = int(fractional_str) / (10 ** len(fractional_str))
            result = ""

            for _ in range(precision):
                fractional *= 2
                if fractional >= 1:
                    result += "1"
                    fractional -= 1
                else:
                    result += "0"
                if fractional == 0:
                    break

            return result if result else "0"
        except (ValueError, Exception):
            return fractional_str


# تست
if __name__ == "__main__":
    converter = BinaryConverter()

    print("=== Programming Language Tests ===")
    tests = [
        "x <= 12",
        "x = 3.14",
        "ali < 3",
        "y != ali",
        "if x == 5",
        "x = a + b * c",
        "a and b or c"
    ]

    for test in tests:
        result = converter.convert_text(test)
        print(f"'{test}': '{result}'")