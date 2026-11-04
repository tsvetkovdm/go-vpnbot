func NewBotAPI(config Config) (*tgbotapi.BotAPI, error) {
	httpClient := &http.Client{}

	if config.ProxyAddr != "" {
		dialer, err := proxy.SOCKS5("tcp", config.ProxyAddr, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("create SOCKS5 dialer: %w", err)
		}

		contextDialer, ok := dialer.(proxy.ContextDialer)
		if !ok {
			return nil, fmt.Errorf("SOCKS5 dialer does not implement ContextDialer")
		}

		httpClient.Transport = &http.Transport{
			DialContext: contextDialer.DialContext,
		}
	}

	bot, err := tgbotapi.NewBotAPIWithClient(config.Token, tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		return nil, fmt.Errorf("create bot API: %w", err)
	}

	return bot, nil
}