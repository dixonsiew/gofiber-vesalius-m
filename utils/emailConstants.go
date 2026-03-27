package utils

const EmailTemplateConstantSignUp = `
    <!DOCTYPE html
    PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
  <html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">
  
  <head>
    <meta charset="UTF-8">
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta name="x-apple-disable-message-reformatting">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta content="telephone=no" name="format-detection">
    <title></title>
    <link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
  </head>
  
  <body>
    <div class="es-wrapper-color">
      <table class="es-wrapper" cellspacing="0" cellpadding="0"
        style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
        <tbody>
          <tr>
            <td class="esd-email-paddings" valign="top">
              <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                        <tbody>
                                          <tr>
                                            <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                                target="_blank"><img class="adapt-img"
                                                  src="https://i.imgur.com/McHoODk.png" alt
                                                  style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                                  ></a></td>
                                          </tr>
                                          <tr>
                                            <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                              <table border="0" width="100%" height="100%" cellpadding="0"
                                                cellspacing="0">
                                                <tbody>
                                                  <tr>
                                                    <td
                                                      style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                    </td>
                                                  </tr>
                                                </tbody>
                                              </table>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20r es-p40l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="540" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text">
                                              <p
                                                style="color: #00585F; line-height: 150%; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 22px; margin-bottom: 0;">
                                                <strong>Verify Your {{email_appname}} Mobile Application Account</strong>
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table class="es-content" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                                <strong>Hey {{first_name}},</strong>
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table class="es-footer" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                Thank you for creating an account with the {{email_appname}} Mobile App. To verify your
                                                account, enter the verification code given below into your {{email_appname}} Mobile App.</p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" class="es-left" align="left">
                                <tbody>
                                  <tr>
                                    <td width="174" class="es-m-p0r es-m-p20b esd-container-frame">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td class="esd-empty-container" style="display: none;"></td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                    <td class="es-hidden" width="20"></td>
                                  </tr>
                                </tbody>
                              </table>
                              <table cellpadding="0" cellspacing="0" class="es-left" align="center">
                                <tbody>
                                  <tr>
                                    <td width="173" class="es-m-p20b esd-container-frame">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td
                                              style="background-color: #337A7F; padding-top: 5px; padding-bottom: 5px; padding-left: 30px; padding-right: 30px;"
                                              class="esd-block-text es-p15t es-p15b">
                                              <p
                                                style="margin-bottom: 20px; margin-top: 20px; letter-spacing: 5px; color: #00474C; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 32px;">
                                                <strong>{{verification_code}}</strong>
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                              <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                                <tbody>
                                  <tr>
                                    <td width="173" class="esd-container-frame">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td class="esd-empty-container" style="display: none;"></td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                <strong>Tips to help protect your password :</strong>
                                              </p>
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                • Never share your password with anyone.<br>• Create passwords that are
                                                hard to guess and don't use personal information. Be sure to include
                                                uppercase and lowercase letters, numbers and symbols.</p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table cellpadding="0" cellspacing="0" class="es-content">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                        width="600">
                        <tbody>
                          <tr>
                            <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                                If you did not create an account with us, please ignore this message.</p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table cellpadding="0" cellspacing="0" class="es-content esd-footer-popover">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                        width="600">
                        <tbody>
                          <tr>
                            <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                Thanks. <br> The {{email_appname}} Mobile App Team</p>
                                            </td>
                                          </tr>
                                          <tr>
                                            <td class="esd-block-spacer es-p20" style="font-size:0">
                                              <table border="0" width="100%" height="100%" cellpadding="0"
                                                cellspacing="0">
                                                <tbody>
                                                  <tr>
                                                    <td
                                                      style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                    </td>
                                                  </tr>
                                                </tbody>
                                              </table>
                                            </td>
                                          </tr>
                                          <tr>
                                            <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                              <tbody>
                                                <tr>
                                                  <td align="center" class="esd-block-text es-p20b"
                                                    esd-links-color="#0A7AFF">
                                                    <p
                                                      style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                      Need help? Reach out to us at <a href="mailto:info@islandhospital.com"
                                                        style="color: #0a7aff;">info@islandhospital.com</a>
                                                    </p>
                                                    <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                      (This is a system-generated email. Please do not reply to this email)
                                                    </p>
                                                  </td>
                                                </tr>
                                              </tbody>
                                            </table>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </body>
  
  </html>
`

const EmailTemplateConstantResetPw = `
    <!DOCTYPE html
    PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
  <html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">
  
  <head>
    <meta charset="UTF-8">
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta name="x-apple-disable-message-reformatting">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta content="telephone=no" name="format-detection">
    <title></title>
    <link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
  </head>
  
  <body>
    <div class="es-wrapper-color">
      <table class="es-wrapper" cellspacing="0" cellpadding="0"
        style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
        <tbody>
          <tr>
            <td class="esd-email-paddings" valign="top">
              <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                        <tbody>
                                          <tr>
                                            <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                                target="_blank"><img class="adapt-img"
                                                  src="https://i.imgur.com/McHoODk.png" alt
                                                  style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                                  ></a></td>
                                          </tr>
                                          <tr>
                                            <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                              <table border="0" width="100%" height="100%" cellpadding="0"
                                                cellspacing="0">
                                                <tbody>
                                                  <tr>
                                                    <td
                                                      style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                    </td>
                                                  </tr>
                                                </tbody>
                                              </table>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20r es-p40l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="540" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text">
                                              <p
                                                style="color: #00585F; line-height: 150%; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 22px; margin-bottom: 0;">
                                                <strong>Reset Your {{email_appname}} Mobile Application Password</strong>
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table class="es-content" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                                <strong>Hey {{first_name}},</strong>
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table class="es-footer" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                We're sending you this email because you requested a password reset.
                                                To continue, enter the verification code given below in {{email_appname}} Mobile App.
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" class="es-left" align="left">
                                <tbody>
                                  <tr>
                                    <td width="174" class="es-m-p0r es-m-p20b esd-container-frame">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td class="esd-empty-container" style="display: none;"></td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                    <td class="es-hidden" width="20"></td>
                                  </tr>
                                </tbody>
                              </table>
                              <table cellpadding="0" cellspacing="0" class="es-right" align="center">
                                <tbody>
                                  <tr>
                                    <td width="173" class="esd-container-frame">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td
                                              style="background-color: #337A7F; padding-top: 5px; padding-bottom: 5px; padding-left: 30px; padding-right: 30px;"
                                              class="esd-block-text es-p15t es-p15b">
                                              <p
                                                style="margin-bottom: 20px; margin-top: 20px; letter-spacing: 5px; color: #00474C; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 32px;">
                                                <strong>{{verification_code}}</strong>
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                <strong>Tips to help protect your password :</strong>
                                              </p>
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                • Never share your password with anyone.<br>• Create passwords that are
                                                hard to guess and don't use personal information. Be sure to include
                                                uppercase and lowercase letters, numbers and symbols.</p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table cellpadding="0" cellspacing="0" class="es-content">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                        width="600">
                        <tbody>
                          <tr>
                            <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                                If you did not request a password reset, please ignore this message.</p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table cellpadding="0" cellspacing="0" class="es-content esd-footer-popover">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                        width="600">
                        <tbody>
                          <tr>
                            <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                Thanks. <br> The {{email_appname}} Mobile App Team</p>
                                            </td>
                                          </tr>
                                          <tr>
                                            <td class="esd-block-spacer es-p20" style="font-size:0">
                                              <table border="0" width="100%" height="100%" cellpadding="0"
                                                cellspacing="0">
                                                <tbody>
                                                  <tr>
                                                    <td
                                                      style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                    </td>
                                                  </tr>
                                                </tbody>
                                              </table>
                                            </td>
                                          </tr>
                                          <tr>
                                            <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                              <tbody>
                                                <tr>
                                                  <td align="center" class="esd-block-text es-p20b"
                                                    esd-links-color="#0A7AFF">
                                                    <p
                                                      style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                      Need help? Reach out to us at <a href="mailto:info@islandhospital.com"
                                                        style="color: #0a7aff;">info@islandhospital.com</a>
                                                    </p>
                                                    <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                      (This is a system-generated email. Please do not reply to this email.)
                                                    </p>
                                                  </td>
                                                </tr>
                                              </tbody>
                                            </table>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </body>
  
  </html>
`

const EmailTemplateConstantLittleKids = `
<!DOCTYPE html
    PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
  <html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">
  
  <head>
    <meta charset="UTF-8">
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta name="x-apple-disable-message-reformatting">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta content="telephone=no" name="format-detection">
    <title></title>
    <link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
  </head>

<body>
  <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
    <strong>Little Explorers Kids Club - Welcome Note</strong>
  </p>
  <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
    Dear {{name}},
  </p>
  <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
    Welcome to the Little Explorers Kids Club!
  </p>
  <br>
  <img src="https://i.imgur.com/3EisZd0.png" alt="LittleExplorersKidsWelcome">
</body>

</html>
`

const EmailTemplateConstantGoldenPearl = `
<!DOCTYPE html
    PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
  <html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">
  
  <head>
    <meta charset="UTF-8">
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta name="x-apple-disable-message-reformatting">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta content="telephone=no" name="format-detection">
    <title></title>
    <link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
  </head>

<body>
  <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
    <strong>Golden Pearl Club - Welcome Note</strong>
  </p>
  <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
    Dear {{name}},
  </p>
  <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
    Welcome to the Golden Pearl Club!
  </p>
  <br>
  <img src="https://i.imgur.com/gTs9p9P.jpg" style="width: 100%; height: 100%" alt="GoldenPearlWelcome">
</body>

</html
`

const EmailTemplateConstantLogisticRequestConfirmation = `
<!DOCTYPE html
    PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
  <html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">
  
  <head>
    <meta charset="UTF-8">
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta name="x-apple-disable-message-reformatting">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta content="telephone=no" name="format-detection">
    <title></title>
    <link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
  </head>
  
  <body>
    <div class="es-wrapper-color">
      <table class="es-wrapper" cellspacing="0" cellpadding="0"
        style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
        <tbody>
          <tr>
            <td class="esd-email-paddings" valign="top">
              <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                        <tbody>
                                          <tr>
                                            <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                                target="_blank"><img class="adapt-img"
                                                  src="https://i.imgur.com/McHoODk.png" alt
                                                  style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                                  ></a></td>
                                          </tr>
                                          <tr>
                                            <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                              <table border="0" width="100%" height="100%" cellpadding="0"
                                                cellspacing="0">
                                                <tbody>
                                                  <tr>
                                                    <td
                                                      style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                    </td>
                                                  </tr>
                                                </tbody>
                                              </table>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table class="es-content" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                                Hi All,
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table class="es-footer" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                We have a new patient submitted airport pickup confirmation in IH Mobile App. Here are the order details : </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                                <tbody>
                                  <tr>
                                    <td width="173" class="esd-container-frame">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td class="esd-empty-container" style="display: none;"></td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                <strong>
                                                  Name: {{requester_name}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  PRN: {{requester_prn}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  Companion: {{with_companion}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  Companion Name: {{companion_name}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  Pickup Date / Time: {{pickup_datetime}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  Request Number: {{logistic_number}}
                                                </strong>
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table cellpadding="0" cellspacing="0" class="es-content">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                        width="600">
                        <tbody>
                          <tr>
                            <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                                Kindly acknowledge the receipt of this request and proceed with processing in our web portal.</p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table cellpadding="0" cellspacing="0" class="es-content esd-footer-popover">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                        width="600">
                        <tbody>
                          <tr>
                            <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                Thanks.</p>
                                            </td>
                                          </tr>
                                          <tr>
                                            <td class="esd-block-spacer es-p20" style="font-size:0">
                                              <table border="0" width="100%" height="100%" cellpadding="0"
                                                cellspacing="0">
                                                <tbody>
                                                  <tr>
                                                    <td
                                                      style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                    </td>
                                                  </tr>
                                                </tbody>
                                              </table>
                                            </td>
                                          </tr>
                                          <tr>
                                            <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                              <tbody>
                                                <tr>
                                                  <td align="center" class="esd-block-text es-p20b"
                                                    esd-links-color="#0A7AFF">
                                                    <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                      (This is a system-generated email. Please do not reply to this email)
                                                    </p>
                                                  </td>
                                                </tr>
                                              </tbody>
                                            </table>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </body>
  
  </html
`

const EmailTemplateConstantLogisticRequestCancellation = `
<!DOCTYPE html
    PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
  <html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">
  
  <head>
    <meta charset="UTF-8">
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta name="x-apple-disable-message-reformatting">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta content="telephone=no" name="format-detection">
    <title></title>
    <link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
  </head>
  
  <body>
    <div class="es-wrapper-color">
      <table class="es-wrapper" cellspacing="0" cellpadding="0"
        style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
        <tbody>
          <tr>
            <td class="esd-email-paddings" valign="top">
              <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                        <tbody>
                                          <tr>
                                            <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                                target="_blank"><img class="adapt-img"
                                                  src="https://i.imgur.com/McHoODk.png" alt
                                                  style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                                  ></a></td>
                                          </tr>
                                          <tr>
                                            <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                              <table border="0" width="100%" height="100%" cellpadding="0"
                                                cellspacing="0">
                                                <tbody>
                                                  <tr>
                                                    <td
                                                      style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                    </td>
                                                  </tr>
                                                </tbody>
                                              </table>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table class="es-content" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                                Hi All,
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table class="es-footer" cellspacing="0" cellpadding="0">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                        style="background-color: #ffffff;">
                        <tbody>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                Please be advised that the patient has cancelled their airport pickup request. Here are the order details :</p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                                <tbody>
                                  <tr>
                                    <td width="173" class="esd-container-frame">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td class="esd-empty-container" style="display: none;"></td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          <tr>
                            <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                <strong>
                                                  Name: {{requester_name}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  PRN: {{requester_prn}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  Companion: {{with_companion}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  Companion Name: {{companion_name}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  Pickup Date / Time: {{pickup_datetime}}
                                                </strong>
                                                <br>
                                                <strong>
                                                  Request Number: {{logistic_number}}
                                                </strong>
                                              </p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table cellpadding="0" cellspacing="0" class="es-content">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                        width="600">
                        <tbody>
                          <tr>
                            <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                                The slot has been release in IH Mobile App. Kindly check and confirm again from your end.</p>
                                            </td>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
              <table cellpadding="0" cellspacing="0" class="es-content esd-footer-popover">
                <tbody>
                  <tr>
                    <td class="esd-stripe">
                      <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                        width="600">
                        <tbody>
                          <tr>
                            <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                              <table cellpadding="0" cellspacing="0" width="100%">
                                <tbody>
                                  <tr>
                                    <td width="560" class="esd-container-frame" valign="top">
                                      <table cellpadding="0" cellspacing="0" width="100%">
                                        <tbody>
                                          <tr>
                                            <td align="left" class="esd-block-text es-p20l">
                                              <p
                                                style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                                Thanks.</p>
                                            </td>
                                          </tr>
                                          <tr>
                                            <td class="esd-block-spacer es-p20" style="font-size:0">
                                              <table border="0" width="100%" height="100%" cellpadding="0"
                                                cellspacing="0">
                                                <tbody>
                                                  <tr>
                                                    <td
                                                      style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                    </td>
                                                  </tr>
                                                </tbody>
                                              </table>
                                            </td>
                                          </tr>
                                          <tr>
                                            <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                              <tbody>
                                                <tr>
                                                  <td align="center" class="esd-block-text es-p20b"
                                                    esd-links-color="#0A7AFF">
                                                    <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                      (This is a system-generated email. Please do not reply to this email)
                                                    </p>
                                                  </td>
                                                </tr>
                                              </tbody>
                                            </table>
                                          </tr>
                                        </tbody>
                                      </table>
                                    </td>
                                  </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </tbody>
              </table>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </body>
  
  </html>
`

const EmailTemplateConstantClubsEventRegistrationToMember = `
<!DOCTYPE html
  PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
  <meta charset="UTF-8">
  <meta content="width=device-width, initial-scale=1" name="viewport">
  <meta name="x-apple-disable-message-reformatting">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta content="telephone=no" name="format-detection">
  <title></title>
  <link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
</head>

<body>
  <div class="es-wrapper-color">
    <table class="es-wrapper" cellspacing="0" cellpadding="0"
      style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
      <tbody>
        <tr>
          <td class="esd-email-paddings" valign="top">
            <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
              <tbody>
                <tr>
                  <td class="esd-stripe">
                    <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                      style="background-color: #ffffff;">
                      <tbody>
                        <tr>
                          <td class="esd-structure es-p20t es-p20r es-p20l">
                            <table cellpadding="0" cellspacing="0" width="100%">
                              <tbody>
                                <tr>
                                  <td width="560" class="esd-container-frame" valign="top">
                                    <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                      <tbody>
                                        <tr>
                                          <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                              target="_blank"><img class="adapt-img"
                                                src="https://i.imgur.com/McHoODk.png" alt
                                                style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                                ></a></td>
                                        </tr>
                                        <tr>
                                          <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                            <table border="0" width="100%" height="100%" cellpadding="0"
                                              cellspacing="0">
                                              <tbody>
                                                <tr>
                                                  <td
                                                    style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                  </td>
                                                </tr>
                                              </tbody>
                                            </table>
                                          </td>
                                        </tr>
                                      </tbody>
                                    </table>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </td>
                </tr>
              </tbody>
            </table>
            <table class="es-content" cellspacing="0" cellpadding="0">
              <tbody>
                <tr>
                  <td class="esd-stripe">
                    <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                      style="background-color: #ffffff;">
                      <tbody>
                        <tr>
                          <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                            <table cellpadding="0" cellspacing="0" width="100%">
                              <tbody>
                                <tr>
                                  <td width="560" class="esd-container-frame" valign="top">
                                    <table cellpadding="0" cellspacing="0" width="100%">
                                      <tbody>
                                        <tr>
                                          <td align="left" class="esd-block-text es-p20l">
                                            <p
                                              style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                              Dear {{member_name}},
                                            </p>
                                          </td>
                                        </tr>
                                      </tbody>
                                    </table>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </td>
                </tr>
              </tbody>
            </table>
            <table class="es-footer" cellspacing="0" cellpadding="0">
              <tbody>
                <tr>
                  <td class="esd-stripe">
                    <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                      style="background-color: #ffffff;">
                      <tbody>
                        <tr>
                          <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                            <table cellpadding="0" cellspacing="0" width="100%">
                              <tbody>
                                <tr>
                                  <td width="560" class="esd-container-frame" valign="top">
                                    <table cellpadding="0" cellspacing="0" width="100%">
                                      <tbody>
                                        <tr>
                                          <td align="left" class="esd-block-text es-p20l">
                                            <p
                                              style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                              We're excited to confirm your registration for the {{event_name}}!</p>
                                          </td>
                                        </tr>
                                      </tbody>
                                    </table>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                          </td>
                        </tr>
                        <tr>
                          <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                            <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                              <tbody>
                                <tr>
                                  <td width="173" class="esd-container-frame">
                                    <table cellpadding="0" cellspacing="0" width="100%">
                                      <tbody>
                                        <tr>
                                          <td class="esd-empty-container" style="display: none;"></td>
                                        </tr>
                                      </tbody>
                                    </table>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </td>
                </tr>
              </tbody>
            </table>
            <table cellpadding="0" cellspacing="0" class="es-content">
              <tbody>
                <tr>
                  <td class="esd-stripe">
                    <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                      width="600">
                      <tbody>
                        <tr>
                          <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                            <table cellpadding="0" cellspacing="0" width="100%">
                              <tbody>
                                <tr>
                                  <td width="560" class="esd-container-frame" valign="top">
                                    <table cellpadding="0" cellspacing="0" width="100%">
                                      <tbody>
                                        <tr>
                                          <td align="left" class="esd-block-text es-p20l">
                                            <p
                                              style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                              If you have any questions, please don't hesitate to reach out to us via email at club@islandhospital.com or by phone at +604 238 3388 (ext 2509).</p>
                                          </td>
                                        </tr>
                                      </tbody>
                                    </table>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </td>
                </tr>
              </tbody>
            </table>
            <table cellpadding="0" cellspacing="0" class="es-content esd-footer-popover">
              <tbody>
                <tr>
                  <td class="esd-stripe">
                    <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                      width="600">
                      <tbody>
                        <tr>
                          <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                            <table cellpadding="0" cellspacing="0" width="100%">
                              <tbody>
                                <tr>
                                  <td width="560" class="esd-container-frame" valign="top">
                                    <table cellpadding="0" cellspacing="0" width="100%">
                                      <tbody>
                                        <tr>
                                          <td align="left" class="esd-block-text es-p20l">
                                            <p
                                              style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                              Thank you for joining us, we look forward to seeing you at the event!</p>
                                          </td>
                                        </tr>
                                        <tr>
                                          <td class="esd-block-spacer es-p20" style="font-size:0">
                                            <table border="0" width="100%" height="100%" cellpadding="0"
                                              cellspacing="0">
                                              <tbody>
                                                <tr>
                                                  <td
                                                    style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                  </td>
                                                </tr>
                                              </tbody>
                                            </table>
                                          </td>
                                        </tr>
                                        <tr>
                                          <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                            <tbody>
                                              <tr>
                                                <td align="center" class="esd-block-text es-p20b"
                                                  esd-links-color="#0A7AFF">
                                                  <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                    (This is a system-generated email. Please do not reply to this email)
                                                  </p>
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </tr>
                                      </tbody>
                                    </table>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </td>
                </tr>
              </tbody>
            </table>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</body>

</html>
`

const EmailTemplateConstantClubsEventRegistrationToIH = `
<!DOCTYPE html
PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
<meta charset="UTF-8">
<meta content="width=device-width, initial-scale=1" name="viewport">
<meta name="x-apple-disable-message-reformatting">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta content="telephone=no" name="format-detection">
<title></title>
<link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
</head>

<body>
<div class="es-wrapper-color">
  <table class="es-wrapper" cellspacing="0" cellpadding="0"
    style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
    <tbody>
      <tr>
        <td class="esd-email-paddings" valign="top">
          <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                    <tbody>
                                      <tr>
                                        <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                            target="_blank"><img class="adapt-img"
                                              src="https://i.imgur.com/McHoODk.png" alt
                                              style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                              ></a></td>
                                      </tr>
                                      <tr>
                                        <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                          <table border="0" width="100%" height="100%" cellpadding="0"
                                            cellspacing="0">
                                            <tbody>
                                              <tr>
                                                <td
                                                  style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-content" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                            Hi All,
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-footer" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            We have a new member registration for {{event_name}}:
                                          </p>
                                          <p style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            Member Name: {{member_name}}
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                            <tbody>
                              <tr>
                                <td width="173" class="esd-container-frame">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td class="esd-empty-container" style="display: none;"></td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table cellpadding="0" cellspacing="0" class="es-content">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                    width="600">
                    <tbody>
                      <tr>
                        <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                            Please get the latest member list event report from the web portal.
                                          </p>
                                          <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                            Thanks.
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
        </td>
      </tr>
    </tbody>
  </table>
</div>
</body>

</html>
`

const EmailTemplateConstantGuestAppointmentConfirmationToIH = `
<!DOCTYPE html
PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
<meta charset="UTF-8">
<meta content="width=device-width, initial-scale=1" name="viewport">
<meta name="x-apple-disable-message-reformatting">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta content="telephone=no" name="format-detection">
<title></title>
<link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
</head>

<body>
<div class="es-wrapper-color">
  <table class="es-wrapper" cellspacing="0" cellpadding="0"
    style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
    <tbody>
      <tr>
        <td class="esd-email-paddings" valign="top">
          <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                    <tbody>
                                      <tr>
                                        <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                            target="_blank"><img class="adapt-img"
                                              src="https://i.imgur.com/McHoODk.png" alt
                                              style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                              ></a></td>
                                      </tr>
                                      <tr>
                                        <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                          <table border="0" width="100%" height="100%" cellpadding="0"
                                            cellspacing="0">
                                            <tbody>
                                              <tr>
                                                <td
                                                  style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-content" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                            Hi All,
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-footer" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            We have a new patient submitted appointment in Island Hospital App. </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                            <tbody>
                              <tr>
                                <td width="173" class="esd-container-frame">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td class="esd-empty-container" style="display: none;"></td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            <strong>
                                              Guest Name: {{guest_name}}
                                            </strong>
                                            <br>
                                            <strong>
                                              Dr Name: {{doctor_name}}
                                            </strong>
                                            <br>
                                            <strong>
                                              Date: {{appointment_date}}
                                            </strong>
                                            <br>
                                            <strong>
                                              Time: {{appointment_time}}
                                            </strong>
                                            <br>
                                            <strong>
                                              Location: {{clinic_location}}
                                            </strong>
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table cellpadding="0" cellspacing="0" class="es-content">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                    width="600">
                    <tbody>
                      <tr>
                        <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                            Kindly get more info from web portal.
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table cellpadding="0" cellspacing="0" class="es-content esd-footer-popover">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                    width="600">
                    <tbody>
                      <tr>
                        <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            Thanks.
                                          </p>
                                        </td>
                                      </tr>
                                      <tr>
                                        <td class="esd-block-spacer es-p20" style="font-size:0">
                                          <table border="0" width="100%" height="100%" cellpadding="0"
                                            cellspacing="0">
                                            <tbody>
                                              <tr>
                                                <td
                                                  style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </td>
                                      </tr>
                                      <tr>
                                        <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                          <tbody>
                                            <tr>
                                              <td align="center" class="esd-block-text es-p20b"
                                                esd-links-color="#0A7AFF">
                                                <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                  (This is a system-generated email. Please do not reply to this email)
                                                </p>
                                              </td>
                                            </tr>
                                          </tbody>
                                        </table>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
        </td>
      </tr>
    </tbody>
  </table>
</div>
</body>

</html>
`

const EmailTemplateConstantPatientFeedbackSubmitted = `
<!DOCTYPE html
PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
<meta charset="UTF-8">
<meta content="width=device-width, initial-scale=1" name="viewport">
<meta name="x-apple-disable-message-reformatting">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta content="telephone=no" name="format-detection">
<title></title>
<link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
</head>

<body>
<div class="es-wrapper-color">
  <table class="es-wrapper" cellspacing="0" cellpadding="0"
    style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
    <tbody>
      <tr>
        <td class="esd-email-paddings" valign="top">
          <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                    <tbody>
                                      <tr>
                                        <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                            target="_blank"><img class="adapt-img"
                                              src="https://i.imgur.com/McHoODk.png" alt
                                              style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                              ></a></td>
                                      </tr>
                                      <tr>
                                        <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                          <table border="0" width="100%" height="100%" cellpadding="0"
                                            cellspacing="0">
                                            <tbody>
                                              <tr>
                                                <td
                                                  style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-content" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                            Hi All,
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-footer" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            We have a new patient feedback submitted in Island Hospital Mobile App.
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                            <tbody>
                              <tr>
                                <td width="173" class="esd-container-frame">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td class="esd-empty-container" style="display: none;"></td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table cellpadding="0" cellspacing="0" class="es-content">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                    width="600">
                    <tbody>
                      <tr>
                        <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                            Please get latest feedback list report from web portal.
                                          </p>
                                          <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                            Thanks.
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
        </td>
      </tr>
    </tbody>
  </table>
</div>
</body>

</html>
`

const EmailTemplateConstantSignupSuccess = `
<!DOCTYPE html
PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
<meta charset="UTF-8">
<meta content="width=device-width, initial-scale=1" name="viewport">
<meta name="x-apple-disable-message-reformatting">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta content="telephone=no" name="format-detection">
<title></title>
<link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
</head>

<body>
<div class="es-wrapper-color">
  <table class="es-wrapper" cellspacing="0" cellpadding="0"
    style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
    <tbody>
      <tr>
        <td class="esd-email-paddings" valign="top">
          <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                    <tbody>
                                      <tr>
                                        <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                            target="_blank"><img class="adapt-img"
                                              src="https://i.imgur.com/McHoODk.png" alt
                                              style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                              ></a></td>
                                      </tr>
                                      <tr>
                                        <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                          <table border="0" width="100%" height="100%" cellpadding="0"
                                            cellspacing="0">
                                            <tbody>
                                              <tr>
                                                <td
                                                  style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-content" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                            Dear {{patient_name}},
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-footer" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            Thanks for signing up, we are pleased to have you join our Island Hospital family.</p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                            <tbody>
                              <tr>
                                <td width="173" class="esd-container-frame">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td class="esd-empty-container" style="display: none;"></td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table cellpadding="0" cellspacing="0" class="es-content">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                    width="600">
                    <tbody>
                      <tr>
                        <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                            Your account has been verified and now active. You can now log in using your {{username}} username.
                                            If you have any question or need further assistance, please don't hesitate to reach out to our customer service at +604-238 3388
                                          </p>
                                          <p
                                            style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 16px;">
                                            Thank you for registering with us. We look forward to serving you.
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table cellpadding="0" cellspacing="0" class="es-content esd-footer-popover">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                    width="600">
                    <tbody>
                      <tr>
                        <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            “Island Hospital , to comfort always.”</p>
                                        </td>
                                      </tr>
                                      <tr>
                                        <td class="esd-block-spacer es-p20" style="font-size:0">
                                          <table border="0" width="100%" height="100%" cellpadding="0"
                                            cellspacing="0">
                                            <tbody>
                                              <tr>
                                                <td
                                                  style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </td>
                                      </tr>
                                      <tr>
                                        <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                          <tbody>
                                            <tr>
                                              <td align="center" class="esd-block-text es-p20b"
                                                esd-links-color="#0A7AFF">
                                                <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                  (This is a system-generated email. Please do not reply to this email)
                                                </p>
                                              </td>
                                            </tr>
                                          </tbody>
                                        </table>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
        </td>
      </tr>
    </tbody>
  </table>
</div>
</body>

</html>
`

const EmailTemplateConstantSuccessOutstandingBillPayment = `
<!DOCTYPE html
PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
<meta charset="UTF-8">
<meta content="width=device-width, initial-scale=1" name="viewport">
<meta name="x-apple-disable-message-reformatting">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta content="telephone=no" name="format-detection">
<title></title>
<link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
</head>

<body>
<div class="es-wrapper-color">
  <table class="es-wrapper" cellspacing="0" cellpadding="0"
    style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
    <tbody>
      <tr>
        <td class="esd-email-paddings" valign="top">
          <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                    <tbody>
                                      <tr>
                                        <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                            target="_blank"><img class="adapt-img"
                                              src="https://i.imgur.com/McHoODk.png" alt
                                              style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"
                                              ></a></td>
                                      </tr>
                                      <tr>
                                        <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                          <table border="0" width="100%" height="100%" cellpadding="0"
                                            cellspacing="0">
                                            <tbody>
                                              <tr>
                                                <td
                                                  style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-content" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                            Dear Valued Customer,
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-footer" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            Thank you for using Island Hospital Mobile App. We are pleased to inform you that {{amount}} online payment via {{payment_method}} is successful for the bill and invoice number {{bill_number}} / {{invoice_number}}.
                                          </p>
                                          <p
                                            style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            Kindly note that this is not an official receipt. Please proceed to the Cashier Counter at Island Hospital to obtain the original official receipt.
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                            <tbody>
                              <tr>
                                <td width="173" class="esd-container-frame">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td class="esd-empty-container" style="display: none;"></td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table cellpadding="0" cellspacing="0" class="es-content esd-footer-popover">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                    width="600">
                    <tbody>
                      <tr>
                        <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                          <tbody>
                                            <tr>
                                              <td align="center" class="esd-block-text es-p20b"
                                                esd-links-color="#0A7AFF">
                                                <p style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                  (This is a system-generated email. Please do not reply to this email)
                                                </p>
                                              </td>
                                            </tr>
                                          </tbody>
                                        </table>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
        </td>
      </tr>
    </tbody>
  </table>
</div>
</body>

</html>
`

const EmailTemplateConstantSuccessPackagePayment = `
<!DOCTYPE html
PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
<meta charset="UTF-8">
<meta content="width=device-width, initial-scale=1" name="viewport">
<meta name="x-apple-disable-message-reformatting">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta content="telephone=no" name="format-detection">
<title></title>
<link href="https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i" rel="stylesheet">
<style>
  .container {
    border: 1px solid #ccc;
    padding: 10px;
  }

  .line {
    display: flex;
    justify-content: space-between;
  }

  .section-order-number {
    width: 100%;
    border: 1px solid #ccc;
    padding: 5px;
    box-sizing: border-box;
  }

  .section-package-expiry-date {
    width: 100%;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
    border-bottom: 1px solid #ccc;
    padding: 5px;
    box-sizing: border-box;
  }

  .section-billing-address {
    width: 100%;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
    border-bottom: 1px solid #ccc;
    padding: 5px;
    box-sizing: border-box;
  }

  .section-left {
    width: 50%;
    border-left: 1px solid #ccc;
    border-bottom: 1px solid #ccc;
    padding: 5px;
    box-sizing: border-box;
  }

  .section-middle {
    width: 25%;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
    border-bottom: 1px solid #ccc;
    padding: 5px;
    box-sizing: border-box;
    text-align: center;
  }

  .section-right {
    width: 25%;
    border-right: 1px solid #ccc;
    border-bottom: 1px solid #ccc;
    padding: 5px;
    box-sizing: border-box;
    text-align: center;
  }

  .section-33 {
    width: 25%;
    border-right: 1px solid #ccc;
    border-bottom: 1px solid #ccc;
    padding: 5px;
    box-sizing: border-box;
    text-align: center;
  }

  .section-66 {
    width: 75%;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
    border-bottom: 1px solid #ccc;
    padding: 5px;
    box-sizing: border-box;
  }
</style>
</head>

<body>
<div class="es-wrapper-color">
  <table class="es-wrapper" cellspacing="0" cellpadding="0"
    style="border-style: solid; border-color: #F3F4F8; padding: 24px; border-width: 20px;">
    <tbody>
      <tr>
        <td class="esd-email-paddings" valign="top">
          <table class="esd-header-popover es-header" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-header-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                    <tbody>
                                      <tr>
                                        <td class="esd-block-image es-p10t es-p10b es-p20l" style="font-size: 0px;"><a
                                            target="_blank"><img class="adapt-img"
                                              src="https://i.imgur.com/McHoODk.png" alt
                                              style="display: block; padding-top: 20px; padding-bottom: 20px; margin-left: auto; margin-right: auto;"></a>
                                        </td>
                                      </tr>
                                      <tr>
                                        <td class="esd-block-spacer es-p20" style="font-size: 0px;">
                                          <table border="0" width="100%" height="100%" cellpadding="0"
                                            cellspacing="0">
                                            <tbody>
                                              <tr>
                                                <td
                                                  style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-content" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-content-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 18px; color: #002e50; margin-bottom: 0;">
                                            Hi {{patient_name}},
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-footer" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td align="left" class="esd-block-text es-p20l">
                                          <p
                                            style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            Thank you for your recent payment. Just to let you know, we've received
                                            your {{order_number}}, and it is now being processed.
                                          </p>
                                          <p
                                            style="color: #002e50; font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 16px;">
                                            Please be reminded to check the redemption period of your purchased
                                            packages and ensure to utilise them before the last date.
                                          </p>
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                            <tbody>
                              <tr>
                                <td width="173" class="esd-container-frame">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td class="esd-empty-container" style="display: none;"></td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <table class="es-footer" cellspacing="0" cellpadding="0">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table class="es-footer-body" width="600" cellspacing="0" cellpadding="0"
                    style="background-color: #ffffff;">
                    <tbody>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <div class="container">
                                          <div class="line">
                                            <div class="section-order-number">
                                              <div
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                {{order_number}} {{date_of_purchase}}
                                              </div>
                                            </div>
                                          </div>
                                          <div class="line">
                                            <div class="section-left"
                                              style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                              <strong>Product</strong>
                                            </div>
                                            <div class="section-middle"
                                              style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                              <strong>Quantity</strong>
                                            </div>
                                            <div class="section-right"
                                              style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                              <strong>Price</strong>
                                            </div>
                                          </div>
                                          <div class="line">
                                            <div class="section-left"
                                              style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                              {{product_name}}
                                            </div>
                                            <div class="section-middle"
                                              style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                              {{product_quantity}}
                                            </div>
                                            <div class="section-right"
                                              style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                              {{product_price}}
                                            </div>
                                          </div>
                                          <div class="line">
                                            <div class="section-66">
                                              <div
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                Subtotal
                                              </div>
                                            </div>
                                            <div class="section-33">
                                              <div
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                {{subtotal_price}}
                                              </div>
                                            </div>
                                          </div>
                                          <div class="line">
                                            <div class="section-66">
                                              <div
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                Payment Method
                                              </div>
                                            </div>
                                            <div class="section-33">
                                              <div
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                {{payment_method}}
                                              </div>
                                            </div>
                                          </div>
                                          <div class="line">
                                            <div class="section-66">
                                              <div
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                Total
                                              </div>
                                            </div>
                                            <div class="section-33">
                                              <div
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                {{total_price}}
                                              </div>
                                            </div>
                                          </div>
                                          <div class="line">
                                            <div class="section-package-expiry-date">
                                              <div
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                <strong>Package Expiry Date: {{package_expiry_date}}</strong>
                                              </div>
                                            </div>
                                          </div>
                                          <div class="line">
                                            <div class="section-billing-address">
                                              <div
                                                style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; font-size: 14px;">
                                                Billing Address: {{billing_address}}
                                              </div>
                                            </div>
                                          </div>
                                        </div>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                      <tr>
                        <td class="esd-structure es-p20t es-p20r es-p20l" align="left">
                          <table cellpadding="0" cellspacing="0" class="es-right" align="right">
                            <tbody>
                              <tr>
                                <td width="173" class="esd-container-frame">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td class="esd-empty-container" style="display: none;"></td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <br>
          <table cellpadding="0" cellspacing="0" class="es-content esd-footer-popover">
            <tbody>
              <tr>
                <td class="esd-stripe">
                  <table style="background-color: #ffffff;" class="es-content-body" cellpadding="0" cellspacing="0"
                    width="600">
                    <tbody>
                      <tr>
                        <td class="es-p20t es-p20r es-p20l esd-structure" align="left">
                          <table cellpadding="0" cellspacing="0" width="100%">
                            <tbody>
                              <tr>
                                <td width="560" class="esd-container-frame" valign="top">
                                  <table cellpadding="0" cellspacing="0" width="100%">
                                    <tbody>
                                      <tr>
                                        <td class="esd-block-spacer es-p20" style="font-size:0">
                                          <table border="0" width="100%" height="100%" cellpadding="0"
                                            cellspacing="0">
                                            <tbody>
                                              <tr>
                                                <td
                                                  style="border-bottom: 1px solid #cccccc; background: unset; height:1px; width:100%; margin:0px 0px 0px 0px;">
                                                </td>
                                              </tr>
                                            </tbody>
                                          </table>
                                        </td>
                                      </tr>
                                      <tr>
                                        <table cellpadding="0" cellspacing="0" width="100%" align="center">
                                          <tbody>
                                            <tr>
                                              <td align="center" class="esd-block-text es-p20b"
                                                esd-links-color="#0A7AFF">
                                                <p
                                                  style="font-family: lato, 'helvetica neue', helvetica, arial, sans-serif; color: #002e50; font-size: 14px; margin-left: auto; margin-right: auto;">
                                                  (This is a system-generated email. Please do not reply to this
                                                  email)
                                                </p>
                                              </td>
                                            </tr>
                                          </tbody>
                                        </table>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
        </td>
      </tr>
    </tbody>
  </table>
</div>
</body>

</html>
`
