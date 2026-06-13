package posts

import "context"

type service struct {
	services []Service
}

func (service *service) Save(
	ctx context.Context,
	post Post,
) error {
	for _, subService := range service.services {
		if subService == nil {
			continue
		}

		if err := subService.Save(ctx, post); err != nil {
			return err
		}
	}

	return nil
}
