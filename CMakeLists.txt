cmake_minimum_required(VERSION 3.21)
project(fcitx5-binaryime VERSION 1.0.0)

find_package(Fcitx5Core REQUIRED)
include("${FCITX_INSTALL_CMAKECONFIG_DIR}/Fcitx5Utils/Fcitx5CompilerSettings.cmake")

install(FILES config/binaryime.conf DESTINATION "${FCITX_INSTALL_PKGDATADIR}/ime-cpp/config")

add_subdirectory(src)