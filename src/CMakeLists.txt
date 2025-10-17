add_library(binaryime SHARED binaryime.cpp)
target_link_libraries(binaryime PRIVATE Fcitx5::Core)
set_target_properties(binaryime PROPERTIES PREFIX "")

install(TARGETS binaryime DESTINATION "${FCITX_INSTALL_LIBDIR}/fcitx5")

# Addon config file
configure_file(binaryime-addon.conf.in binaryime-addon.conf)
install(FILES "${CMAKE_CURRENT_BINARY_DIR}/binaryime-addon.conf" 
        RENAME binaryime.conf 
        DESTINATION "${FCITX_INSTALL_PKGDATADIR}/addon")

# Input Method registration file
install(FILES "binaryime.conf" 
        DESTINATION "${FCITX_INSTALL_PKGDATADIR}/inputmethod")