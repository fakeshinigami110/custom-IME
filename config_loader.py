import configparser
import os


def load_engine_config():
    """لود تنظیمات اصلی"""
    config = configparser.ConfigParser()
    config_path = './binary_ime.conf'

    if os.path.exists(config_path):
        config.read(config_path)
    else:
        config['Engine'] = {
            'convert_numbers': 'true',
            'unknown_chars_behavior': 'keep',
            'add_spaces': 'true'
        }
        config['Mapping'] = {
            'mapping_file': './mappings/default.conf',
            'keywords_file': './mappings/keywords.conf'
        }
        with open(config_path, 'w') as f:
            config.write(f)

    return config


def load_mapping_config():
    """لود مپینگ کاراکترها"""
    engine_config = load_engine_config()
    mapping_file = engine_config.get('Mapping', 'mapping_file')

    mapping = {}
    if os.path.exists(mapping_file):
        with open(mapping_file, 'r') as f:
            for line in f:
                line = line.strip()
                if line and '=' in line and not line.startswith('#'):
                    key, value = line.split('=', 1)
                    key = key.strip()
                    value = value.strip()
                    mapping[key] = value
    else:
        mapping = create_default_mapping()

    return mapping


def load_keywords_config():
    """لود کلمات کلیدی و عملگرها"""
    engine_config = load_engine_config()
    keywords_file = engine_config.get('Mapping', 'keywords_file', fallback='./mappings/keywords.conf')

    keywords = {}
    operators = {}

    if os.path.exists(keywords_file):
        with open(keywords_file, 'r', encoding='utf-8') as f:
            content = f.read()

        # پارس دستی برای جلوگیری از مشکل با =
        current_section = None
        for line in content.split('\n'):
            line = line.strip()
            if not line or line.startswith('#'):
                continue

            if line.startswith('[') and line.endswith(']'):
                current_section = line[1:-1].strip()
            elif current_section == 'Keywords' and '=' in line:
                key, value = line.split('=', 1)
                key = key.strip()
                value = value.strip()
                keywords[key] = value
            elif current_section == 'Operators' and ':' in line:
                key, value = line.split(':', 1)
                key = key.strip()
                value = value.strip()
                operators[key] = value

    return keywords, operators


def create_default_mapping():
    """ایجاد مپینگ پیش‌فرض"""
    base_chars = "abcdefghijklmnopqrstuvwxyz.,?!':"
    mapping = {}

    for i, char in enumerate(base_chars):
        mapping[char] = format(i, '05b')

    return mapping